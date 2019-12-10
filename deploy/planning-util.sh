#!/usr/bin/env bash

#################################################################################################################
#
#   This script is created for demo purpose.
#   Usage:
#
#################################################################################################################

show_usage()
{
    cat << __EOF__

    Usage:
        Requirements:
            --namespace <space> target namespace name [ex: --namespace nginx]
            # we currently support Deployment name, Deployment Configuration name, or StatefulSet name.
            --pod-name <space> target pod name [ex: --pod-name nginx-stable-1-lvr4d]
        Operations:
            --get-current-pod-resources
            --get-pod-planning
            --generate-patch
            --apply-patch <space> patch file full path [ex: --apply-patch /tmp/planning-util/patch.yml]

__EOF__
    exit 1
}

is_pod_ready()
{
  [[ "$(kubectl get po "$1" -n "$2" -o 'jsonpath={.status.conditions[?(@.type=="Ready")].status}')" == 'True' ]]
}

pods_ready()
{
  [[ "$#" == 0 ]] && return 0

  namespace="$1"

  kubectl get pod -n $namespace \
    -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.conditions[?(@.type=="Ready")].status}{"\n"}{end}' \
      | while read name status _junk; do
          if [ "$status" != "True" ]; then
            echo "Waiting pod $name in namespace $namespace to be ready..."
            return 1
          fi
        done || return 1

  return 0
}

leave_prog()
{
    if [ ! -z "$(ls -A $file_folder)" ]; then      
        echo -e "\n$(tput setaf 6)Downloaded YAML files are located under $file_folder $(tput sgr 0)"
    fi
 
    cd $current_location > /dev/null
}

check_version()
{
    openshift_required_minor_version="9"
    k8s_required_version="11"

    oc version 2>/dev/null|grep "oc v"|grep -q " v[4-9]"
    if [ "$?" = "0" ];then
        # oc version is 4-9, passed
        return 0
    fi

    oc version 2>/dev/null|grep "oc v"|grep -q " v[0-2]"
    if [ "$?" = "0" ];then
        # oc version is 0-2, failed
        echo -e "\n$(tput setaf 1)Error! OpenShift version less than 3.$openshift_required_minor_version is not supported by Federator.ai$(tput sgr 0)"
        exit 5
    fi

    # oc major version = 3
    openshift_minor_version=`oc version 2>/dev/null|grep "oc v"|cut -d '.' -f2`
    # k8s version = 1.x
    k8s_version=`kubectl version 2>/dev/null|grep Server|grep -o "Minor:\"[0-9]*\""|cut -d '"' -f2`

    if [ "$openshift_minor_version" != "" ] && [ "$openshift_minor_version" -lt "$openshift_required_minor_version" ]; then
        echo -e "\n$(tput setaf 1)Error! OpenShift version less than 3.$openshift_required_minor_version is not supported by Federator.ai$(tput sgr 0)"
        exit 5
    elif [ "$openshift_minor_version" = "" ] && [ "$k8s_version" != "" ] && [ "$k8s_version" -lt "$k8s_required_version" ]; then
        echo -e "\n$(tput setaf 1)Error! Kubernetes version less than 1.$k8s_required_version is not supported by Federator.ai$(tput sgr 0)"
        exit 6
    fi
}


wait_until_pods_ready()
{
  period="$1"
  interval="$2"
  namespace="$3"
  target_pod_number="$4"

  wait_pod_creating=1
  for ((i=0; i<$period; i+=$interval)); do

    if [[ "$wait_pod_creating" = "1" ]]; then
        # check if pods created
        if [[ "`kubectl get po -n $namespace 2>/dev/null|wc -l`" -ge "$target_pod_number" ]]; then
            wait_pod_creating=0
            echo -e "\nChecking pods..."
        else
            echo "Waiting for pods in namespace $namespace to be created..."
        fi
    else
        # check if pods running
        if pods_ready $namespace; then
            echo -e "\nAll $namespace pods are ready."
            return 0
        fi
        echo "Waiting for pods in namespace $namespace to be ready..."
    fi

    sleep "$interval"
    
  done

  echo -e "\n$(tput setaf 1)Warning!! Waited for $period seconds, but all pods are not ready yet. Please check $namespace namespace$(tput sgr 0)"
  leave_prog
  exit 4
}

get_k8s_rest_api_node_port()
{
    #K8S
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get REST API service NodePort...$(tput sgr 0)"
    https_node_port="`kubectl get svc -n $install_namespace |grep -E -o "5056:.{0,22}"|cut -d '/' -f1|cut -d ':' -f2`"
    if [ "$https_node_port" = "" ]; then
        echo -e "\n$(tput setaf 1)Error! Can't find NodePort of REST API service.$(tput sgr 0)"
        leave_prog
        exit 8 
        
    fi
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration get_k8s_rest_api_node_port = $duration" >> $debug_log
}

check_rest_api_url()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Checking REST API URL...$(tput sgr 0)"
    if [ "$openshift_minor_version" != "" ]; then
        # OpenShift
        api_url="`oc get route -n $install_namespace | grep "federatorai-rest"|awk '{print $2}'`"
        if [ "$api_url" = "" ]; then
            echo -e "\n$(tput setaf 1)Error! Can't get REST API URL.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        api_url="https://$api_url"
    else
        # K8S
        get_k8s_rest_api_node_port
        read -r -p "$(tput setaf 2)Please input REST API service external IP: $(tput sgr 0) " rest_ip </dev/tty
        if [ "$rest_ip" != "" ]; then
            api_url="https://$rest_ip:$https_node_port"
            echo "api_url = $api_url"
        else
            echo -e "\n$(tput setaf 1)Error! Please input correct REST API service IP.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    fi
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration check_rest_api_url = $duration" >> $debug_log
    echo "REST API URL = $api_url"
}

rest_api_login()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Login to REST API...$(tput sgr 0)"
    check_api_url
    #echo "curl -sS -k -X POST \"$api_url/apis/v1/users/login\" -H \"accept: application/json\" -H \"authorization: Basic YWRtaW46YWRtaW4=\" |jq '.accessToken'|tr -d \"\"\""
    access_token=`curl -sS -k -X POST "$api_url/apis/v1/users/login" -H "accept: application/json" -H "authorization: Basic YWRtaW46YWRtaW4=" |jq '.accessToken'|tr -d "\""`
    check_user_token

    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration rest_api_login = $duration" >> $debug_log
}

check_api_url()
{
    if [ "$api_url" = "" ]; then
        echo -e "\n$(tput setaf 1)Error! REST API URL is empty.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
}

check_user_token()
{
    if [ "$access_token" = "" ]; then
        echo -e "\n$(tput setaf 1)Error! User token is empty.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
}

check_cluster_name()
{
    if [ "$cluster_name" = "" ]; then
        echo -e "\n$(tput setaf 1)Error! cluster name is empty.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
}

rest_api_get_cluster_name()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get cluster name...$(tput sgr 0)"
    cluster_name=`curl -sS -k -X GET "$api_url/apis/v1/resources/clusters" -H "accept: application/json" -H "Authorization: Bearer $access_token" |jq '.data[].name'|tr -d "\""`
    check_cluster_name

    echo "cluster_name = $cluster_name"
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration rest_api_get_cluster_name = $duration" >> $debug_log
}

rest_api_get_pod_planning()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get planning for pod ($target_pod_name) in ns ($target_namespace)...$(tput sgr 0)"
    interval_start_time="$start"
    interval_end_time=$(($interval_start_time + 3599)) #59min59sec
    granularity="3600"
    type="recommendation"

    planning_values=`curl -sS -k -X GET "$api_url/apis/v1/plannings/clusters/$cluster_name/namespaces/$target_namespace/pods?granularity=$granularity&type=$type&names=$target_pod_name&limit=1&order=desc&startTime=$interval_start_time&endTime=$interval_end_time" -H "accept: application/json" -H "Authorization: Bearer $access_token" |jq '.plannings[].containerPlannings[0]|"\(.limitPlannings.CPU_USAGE_SECONDS_PERCENTAGE[].numValue) \(.requestPlannings.CPU_USAGE_SECONDS_PERCENTAGE[].numValue) \(.limitPlannings.MEMORY_USAGE_BYTES[].numValue) \(.requestPlannings.MEMORY_USAGE_BYTES[].numValue)"'|tr -d "\""`
    limits_cpu="`echo $planning_values |awk '{print $1}'`"
    requests_cpu="`echo $planning_values |awk '{print $2}'`"
    limits_memory="`echo $planning_values |awk '{print $3}'`"
    requests_memory="`echo $planning_values |awk '{print $4}'`"
    echo "Planning: ----------------------------------"
    echo "resources.limits.cpu = $limits_cpu(m)"
    echo "resources.limits.momory = $limits_memory(byte)"
    echo "resources.requests.cpu = $requests_cpu(m)"
    echo "resources.requests.memory = $requests_memory(byte)"
    echo "--------------------------------------------"

    if [ "$limits_cpu" = "" ] || [ "$requests_cpu" = "" ] || [ "$limits_memory" = "" ] || [ "$requests_memory" = "" ]; then
        echo -e "\n$(tput setaf 1)Error! Failed to get pod ($target_pod_name) planning. Missing value.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
    
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration rest_api_get_pod_planning = $duration" >> $debug_log
}

get_needed_info()
{
    check_rest_api_url
    rest_api_login
    rest_api_get_cluster_name
    get_controller_info
}

get_owner_reference()
{
    local kind="$1"
    local name="$2"
    local owner_ref=`kubectl get $kind $name -n $target_namespace -o json | jq -r '.metadata.ownerReferences[] | "\(.controller) \(.kind) \(.name)"' 2>/dev/null`
    echo "$owner_ref"
}

get_controller_info()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get pod controller type and name...$(tput sgr 0)"
    owner_kind="pod"
    owner_name="$target_pod_name"
    fist_run="y"
    while true
    do
        owner_ref=$(get_owner_reference $owner_kind $owner_name)
        if [ "$fist_run" = "y" ] && [ "$owner_ref" = "" ]; then
            # Pod # First run
            echo -e "\n$(tput setaf 1)Error! Can't find pod ($target_pod_name) ownerReferences in namespace $target_namespace$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        fist_run="n"
        if [ "$owner_ref" != "" ]; then
            owner_kind="`echo $owner_ref |grep 'true'|awk '{print $2}'`"
            owner_name="`echo $owner_ref |grep 'true'|awk '{print $3}'`"
            if [ "$owner_kind" = "DeploymentConfig" ] || [ "$owner_kind" = "Deployment" ] || [ "$owner_kind" = "StatefulSet" ]; then
                break
            fi
        else
            break
        fi
    done

    echo "target_namespace = $target_namespace"
    echo "target_pod_name = $target_pod_name"
    echo "owner_reference_kind = $owner_kind"
    echo "owner_reference_name = $owner_name"

    if [ "$owner_kind" != "DeploymentConfig" ] && [ "$owner_kind" != "Deployment" ] && [ "$owner_kind" != "StatefulSet" ]; then
        echo -e "\n$(tput setaf 1)Error! Only support DeploymentConfig, Deployment, or StatefulSet for now.$(tput sgr 0)"
        leave_prog
        exit 8
    fi

    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration get_controller_info = $duration" >> $debug_log
}

display_pod_resources()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get current pod resource settings...$(tput sgr 0)"
    echo "target_namespace= $target_namespace"
    echo "target_pod_name= $target_pod_name"
    echo "--------------------------------------------"
    kubectl get pod $target_pod_name -n $target_namespace -o json |jq '.spec.containers[].resources'
    echo "--------------------------------------------"
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration display_pod_resources = $duration" >> $debug_log
}

while getopts "h-:" o; do
    case "${o}" in
        -)
            case "${OPTARG}" in
                namespace)
                    target_namespace="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    if [ "$target_namespace" = "" ]; then
                        echo -e "\n$(tput setaf 1)Error! Missing --${OPTARG} value$(tput sgr 0)"
                        show_usage
                        exit
                    fi
                    ;;
                pod-name)
                    target_pod_name="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    if [ "$target_pod_name" = "" ]; then
                        echo -e "\n$(tput setaf 1)Error! Missing --${OPTARG} value$(tput sgr 0)"
                        show_usage
                        exit
                    fi
                    ;;
                get-current-pod-resources)
                    get_current_pod_resources="y"
                    ;;
                get-pod-planning)
                    get_pod_planning="y"
                    ;;
                generate-patch)
                    generate_patch="y"
                    ;;
                apply-patch)
                    apply_patch="y"
                    patch_path="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    if [ "$patch_path" = "" ]; then
                        echo -e "\n$(tput setaf 1)Error! Missing --${OPTARG} value$(tput sgr 0)"
                        show_usage
                        exit
                    fi
                    ;;
                *)
                    echo -e "\n$(tput setaf 1)Error! Unknown option --${OPTARG}$(tput sgr 0)"
                    exit
                    ;;
            esac;;
        h)
            show_usage
            ;;
        *)
            echo -e "\n$(tput setaf 1)Error! wrong paramter.$(tput sgr 0)"
            exit 5
            ;;
    esac
done

[ "$target_namespace" = "" ] && show_usage
[ "$target_pod_name" = "" ] && show_usage

if [ "$get_current_pod_resources" = "" ] && [ "$get_pod_planning" = "" ] && [ "$generate_patch" = "" ] && [ "$apply_patch" = "" ]; then
    echo -e "\n$(tput setaf 1)Error! At least one operation must be specified.$(tput sgr 0)"
    show_usage
fi

[ "$get_current_pod_resources" = "" ] && get_current_pod_resources="n"
[ "$get_pod_planning" = "" ] && get_pod_planning="n"
[ "$generate_patch" = "" ] && generate_patch="n"
[ "$apply_patch" = "" ] && apply_patch="n"

echo "target_namespace = $target_namespace"
echo "target_pod_name = $target_pod_name"
echo "get_current_pod_resources = $get_current_pod_resources"
echo "get_pod_planning = $get_pod_planning"
echo "generate_patch =$generate_patch"
echo "apply_patch =$apply_patch"
echo "patch_path =$patch_path"

kubectl version|grep -q "^Server"
if [ "$?" != "0" ];then
    echo -e "\nPlease login to kubernetes first."
    exit
fi

which curl > /dev/null 2>&1
if [ "$?" != "0" ];then
    echo -e "\n$(tput setaf 1)Abort, \"curl\" command is needed for this tool.$(tput sgr 0)"
    exit
fi

which jq > /dev/null 2>&1
if [ "$?" != "0" ];then
    echo -e "\n$(tput setaf 1)Abort, \"jq\" command is needed for this tool.$(tput sgr 0)"
    echo "You may issue following commands to install jq."
    echo "1. wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 -O jq"
    echo "2. chmod +x jq"
    echo "3. mv jq /usr/local/bin"
    echo "4. rerun the script"
    exit
fi

echo "Checking environment version..."
check_version
echo "...Passed"

install_namespace="`kubectl get pods --all-namespaces |grep "alameda-ai-"|awk '{print $1}'|head -1`"

if [ "$install_namespace" = "" ];then
    echo -e "\n$(tput setaf 1)Error! Please Install Federatorai before running this script.$(tput sgr 0)"
    exit 3
fi

file_folder="/tmp/recommendation_retriever"
debug_log="debug.log"

rm -rf $file_folder
mkdir -p $file_folder
current_location=`pwd`
cd $file_folder
echo "Receiving command '$0 $@'" >> $debug_log

get_needed_info

if [ "$get_current_pod_resources" = "y" ];then  
    display_pod_resources
fi

if [ "$get_pod_planning" = "y" ];then
    rest_api_get_pod_planning
fi

if [ "$generate_patch" = "y" ];then
    echo ""
fi

if [ "$apply_patch" = "y" ];then
    echo ""
fi


leave_prog
exit 0
