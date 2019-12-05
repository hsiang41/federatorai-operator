#!/usr/bin/env bash

#################################################################################################################
#
#   This script is created for demo purpose.
#   Usage:
#
#################################################################################################################

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

enable_rest_api_node_port()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Enabling REST API service NodePort if needed ...$(tput sgr 0)"
    https_node_port="`kubectl get svc -n $install_namespace |grep -E -o "5056:.{0,22}"|cut -d '/' -f1|cut -d ':' -f2`"
    if [ "$https_node_port" = "" ]; then
        # K8S
        # No Node Port found for 5056 (https)
        node_port_patch="rest_api_node_port.yaml"
        cat > ${node_port_patch} << __EOF__
spec:
  serviceExposures:
  - name: federatorai-rest
    nodePort:
      ports:
      - nodePort: 30055
        port: 5055
      - nodePort: 30056
        port: 5056
    type: NodePort
__EOF__
        alamedaservice_name="`kubectl get alamedaservice -n $install_namespace -o name`"
        if [ "$alamedaservice_name" = "" ]; then
            echo -e "\n$(tput setaf 1)Error! Can't get alamedaservice name.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        kubectl patch -n $install_namespace $alamedaservice_name --type merge --patch "$(cat ${node_port_patch})"
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! Patch alamedaservice failed.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        echo "Checking..."
        sleep 20
        https_node_port="`kubectl get svc -n $install_namespace |grep -E -o "5056:.{0,22}"|cut -d '/' -f1|cut -d ':' -f2`"
        if [ "$https_node_port" = "" ]; then
            echo -e "\n$(tput setaf 1)Error! Still can't find NodePort of REST API service.$(tput sgr 0)"
            leave_prog
            exit 8 
        fi
    fi
    echo "Done."
    end=`date +%s`
    duration=$((end-start))
    echo "Duration enable_rest_api_node_port = $duration" >> $debug_log
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
        enable_rest_api_node_port
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
    check_api_url
    check_user_token
    cluster_name=`curl -sS -k -X GET "$api_url/apis/v1/resources/clusters" -H "accept: application/json" -H "Authorization: Bearer $access_token" |jq '.data[].name'|tr -d "\""`
    check_cluster_name
    
    echo "Done."
    echo "cluster_name = $cluster_name"
    end=`date +%s`
    duration=$((end-start))
    echo "Duration rest_api_get_cluster_name = $duration" >> $debug_log
}

rest_api_get_pod_planning()
{
    start=`date +%s`
    echo -e "\n$(tput setaf 6)Get nginx pod planning...$(tput sgr 0)"
    check_api_url
    check_user_token
    check_cluster_name
    nginx_pod_name=`kubectl get pods -n $nginx_namespace -o name|head -1|cut -d '/' -f2`
    start_time=`date +%s`
    end_time=$(($start_time + 3599)) #59min59sec
    echo "script start_time = $start_time"
    echo "script end_time = $end_time"
    granularity="3600"
    type="recommendation"
    curl -sS -k -X GET "$api_url/apis/v1/plannings/clusters/$cluster_name/namespaces/$nginx_namespace/pods?granularity=$granularity&type=$type&names=$nginx_pod_name&limit=10&order=desc&startTime=$start_time&endTime=$end_time" -H "accept: application/json" -H "Authorization: Bearer $access_token"| python -mjson.tool

    echo "Done."
    echo "cluster_name = $cluster_name"
    end=`date +%s`
    duration=$((end-start))
    echo "Duration rest_api_get_cluster_name = $duration" >> $debug_log
}


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
nginx_namespace="nginx-preloader-sample"
debug_log="debug.log"

rm -rf $file_folder
mkdir -p $file_folder
current_location=`pwd`
cd $file_folder
echo "Receiving command '$0 $@'" >> $debug_log

check_rest_api_url
rest_api_login
rest_api_get_cluster_name
rest_api_get_pod_planning

leave_prog
exit 0
