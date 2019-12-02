#!/bin/bash

#################################################################################################################
#
#   This script is created for demo purpose.
#   Usage:
#       [-p] # Prepare environment
#       [-c] # clean environment for preloader test
#       [-e] # Enable preloader pod
#       [-r] # Run preloader
#       [-f future data point (hour)] # Run preloader future mode
#       [-d] # Disable & Remove preloader
#       [-v] # Revert environment to normal mode
#       [-h] # Display script usage
#
#################################################################################################################

show_usage()
{
    cat << __EOF__

    Usage:
        [-p] # Prepare environment
        [-c] # clean environment for preloader test
        [-e] # Enable preloader pod
        [-r] # Run preloader
        [-f future data point (hour)] # Run preloader future mode
        [-d] # Disable & Remove preloader
        [-v] # Revert environment to normal mode
        [-h] # Display script usage

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

wait_until_data_pump_finish()
{
  period="$1"
  interval="$2"
  type="$3"

  for ((i=0; i<$period; i+=$interval)); do
    if [ "$type" = "future" ]; then
        echo "Waiting for data pump (future mode) to finish ..."
        kubectl logs -n $install_namespace $current_preloader_pod_name | grep -q "Completed to loader container future metrics data"
        if [ "$?" = "0" ]; then
            echo -e "\n$(tput setaf 6)Data pump (future mode) is finish.$(tput sgr 0)"
            return 0
        fi
    else #historical mode
        echo "Waiting for data pump to finish ..."
        kubectl logs -n $install_namespace $current_preloader_pod_name | grep -q "Succeed to generate pods historical metrics"
        if [ "$?" = "0" ]; then
            echo -e "\n$(tput setaf 6)Data pump is finish.$(tput sgr 0)"
            return 0
        fi
    fi
    
    sleep "$interval"
  done

  echo -e "\n$(tput setaf 1)Warning!! Waited for $period seconds, but data pump is still running.$(tput sgr 0)"
  leave_prog
  exit 4
}

get_current_preloader_name()
{
    current_preloader_pod_name=""
    current_preloader_pod_name="`kubectl get pods -n $install_namespace |grep "federatorai-agent-preloader-"|awk '{print $1}'|head -1`"
    echo "current_preloader_pod_name = $current_preloader_pod_name"
}

delete_all_alamedascaler()
{
    echo -e "\n$(tput setaf 6)Deleting old alamedascaler if necessary...$(tput sgr 0)"
    while read alamedascaler_name alamedascaler_ns
    do
        if [ "$alamedascaler_name" = "" ] || [ "$alamedascaler_ns" = "" ]; then
           continue
        fi

        kubectl delete alamedascaler $alamedascaler_name -n $alamedascaler_ns
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in deleting old alamedascaler named $alamedascaler_name in ns $alamedascaler_ns.$(tput sgr 0)"
            exit 8
        fi
    done <<< "$(kubectl get alamedascaler --all-namespaces --output jsonpath='{range .items[*]}{"\n"}{.metadata.name}{"\t"}{.metadata.namespace}' 2>/dev/null)"
    echo "Done"
}

run_preloader_command()
{
    echo -e "\n$(tput setaf 6)Run preloader...$(tput sgr 0)"
    get_current_preloader_name
    if [ "$current_preloader_pod_name" = "" ]; then
        echo -e "\n$(tput setaf 1)ERROR! Can't find installed preloader pod.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
    kubectl exec -n $install_namespace $current_preloader_pod_name -- /opt/alameda/federatorai-agent/bin/transmitter enable
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Error in executing preloader enable command.$(tput sgr 0)"
        exit 8
    fi
    echo "Checking..."
    sleep 20
    kubectl logs -n $install_namespace $current_preloader_pod_name | grep -i "Start PreLoader agent"
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Preloader pod is not running correctly. Please contact support stuff$(tput sgr 0)"
        leave_prog
        exit 5
    fi

    wait_until_data_pump_finish 1200 60 "historical"
    echo "Done."
}

run_futuremode_preloader()
{
    echo -e "\n$(tput setaf 6)Run future mode preloader...$(tput sgr 0)"
    get_current_preloader_name
    if [ "$current_preloader_pod_name" = "" ]; then
        echo -e "\n$(tput setaf 1)ERROR! Can't find installed preloader pod.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
    
    kubectl exec -n $install_namespace $current_preloader_pod_name -- /opt/alameda/federatorai-agent/bin/transmitter loadfuture --hours=$future_mode_length
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Error in executing preloader loadfuture command.$(tput sgr 0)"
        exit 8
    fi

    echo "Checking..."
    sleep 10
    wait_until_data_pump_finish 1200 60 "future"
    echo "Done."
}

reschedule_dispatcher()
{
    echo -e "\n$(tput setaf 6)Reschedule alameda-ai dispatcher...$(tput sgr 0)"
    current_dispatcher_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-ai-dispatcher-"|awk '{print $1}'|head -1`"
    if [ "$current_dispatcher_pod_name" = "" ]; then
        echo -e "\n$(tput setaf 1)ERROR! Can't find alameda-ai dispatcher pod.$(tput sgr 0)"
        leave_prog
        exit 8
    fi

    kubectl delete pod -n $install_namespace $current_dispatcher_pod_name
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Error in deleting dispatcher pod.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
    echo ""
    wait_until_pods_ready 600 30 $install_namespace 5
    echo "Done."

}

get_grafana_route()
{
    if [ "$openshift_minor_version" != "" ] ; then
        link=`oc get route -n $1 2>/dev/null|grep grafana|awk '{print $2}'`
        if [ "$link" != "" ] ; then
            echo -e "\n========================================"
            echo -e "\n$(tput setaf 2)Great! Prediction/Planning jobs are triggered. You could grab a cup of coffee."
            echo "Access GUI through $(tput setaf 6)http://${link} $(tput sgr 0)"
            echo "Default login credential is $(tput setaf 6)admin/admin$(tput sgr 0)"
            echo "========================================"
        else
            echo "Warning! Failed to obtain grafana route address."
        fi
    fi
}

patch_datahub_for_preloader()
{
    echo -e "\n$(tput setaf 6)Starting patch datahub for preloader...$(tput sgr 0)"
    kubectl get deployment/alameda-datahub -n $install_namespace -o yaml|grep ALAMEDA_DATAHUB_APIS_METRICS_SOURCE -A1|grep -q influxdb
    if [ "$?" != "0" ]; then
        kubectl set env deployment/alameda-datahub -n $install_namespace ALAMEDA_DATAHUB_APIS_METRICS_SOURCE=influxdb
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in patching datahub pod.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        echo ""
        wait_until_pods_ready 600 30 $install_namespace 5
    fi
    echo "Done"
}

patch_datahub_back_to_normal()
{
    echo -e "\n$(tput setaf 6)Starting roll back datahub...$(tput sgr 0)"
    kubectl get deployment/alameda-datahub -n $install_namespace -o yaml|grep ALAMEDA_DATAHUB_APIS_METRICS_SOURCE -A1|grep -q prometheus
    if [ "$?" != "0" ]; then
        kubectl set env deployment/alameda-datahub -n $install_namespace ALAMEDA_DATAHUB_APIS_METRICS_SOURCE=prometheus
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in rolling back datahub pod.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        echo ""
        wait_until_pods_ready 600 30 $install_namespace 5
    fi
    echo "Done"
}

check_influxdb_retention()
{
    echo -e "\n$(tput setaf 6)Starting check retention policy...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "show retention policies"|grep "autogen"|grep -q "3600h"
    if [ "$?" != "0" ]; then
        echo -e "\n$(tput setaf 1)Error! retention policy of alameda_metric pod is not 3600h.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
    echo "Done"
}

patch_grafana_for_preloader(){
    echo -e "\n$(tput setaf 6)Starting add flag for grafana ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "select * from grafana_config order by time desc limit 1" 2>/dev/null|grep -q true
    if [ "$?" != "0" ]; then
        kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -execute "show databases" |grep -q "alameda_metric"
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! Can't find alameda_metric in influxdb.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "insert grafana_config preloader=true"
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! add flag for grafana is failed.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    fi
    echo "Done"
}

patch_grafana_back_to_normal(){
    echo -e "\n$(tput setaf 6)Starting add flag to roll back grafana ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "select * from grafana_config order by time desc limit 1" 2>/dev/null|grep -q false
    if [ "$?" != "0" ]; then
        kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -execute "show databases" |grep -q "alameda_metric"
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! Can't find alameda_metric in influxdb.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "insert grafana_config preloader=false"
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! add flag to roll back grafana is failed.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    fi
    echo "Done"
}

verify_metrics_exist()
{
    echo -e "\n$(tput setaf 6)Starting verify metrics in influxdb ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    metrics_list=$(kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "show measurements")
    metrics_num=$(echo "$metrics_list"| egrep "container_cpu|container_memory|node_cpu|node_memory" |wc -l)
    echo "metrics_num = $metrics_num"
    if [ "$metrics_num" -lt "4" ]; then
        echo -e "\n$(tput setaf 1)Error! metrics in alameda_metric is not complete.$(tput sgr 0)"
        echo "$metrics_list"
        exit 8
    fi
    echo "Done"
}

delete_nginx_example()
{
    echo -e "\n$(tput setaf 6)Deleting nginx sample ...$(tput sgr 0)"
    dc_name="`kubectl get dc -n $nginx_ns 2>/dev/null|grep -v "NAME"|awk '{print $1}'`"
    if [ "$dc_name" != "" ]; then
        kubectl delete dc $dc_name -n $nginx_ns
    fi
    deploy_name="`kubectl get deploy -n $nginx_ns 2>/dev/null|grep -v "NAME"|awk '{print $1}'`"
    if [ "$deploy_name" != "" ]; then
        kubectl delete deploy $deploy_name -n $nginx_ns
    fi
    kubectl get ns $nginx_ns >/dev/null 2>&1
    if [ "$?" = "0" ]; then
        kubectl delete ns $nginx_ns
    fi
    echo "Done"
}

new_nginx_example()
{
    echo -e "\n$(tput setaf 6)Creating new nginx sample pod ...$(tput sgr 0)"

    if [[ "`kubectl get po -n $nginx_ns 2>/dev/null|grep -v "NAME"|grep "Running"|wc -l`" -gt "0" ]]; then
        echo "nginx-preloader-sample namespace and pod are already exist."
    else
        if [ "$openshift_minor_version" != "" ]; then
            # OpenShift
            oc new-project $nginx_ns
            oc new-app twalter/openshift-nginx:stable --name nginx-stable
            if [ "$?" != "0" ]; then
                echo -e "\n$(tput setaf 1)Error! create nginx app failed.$(tput sgr 0)"
                leave_prog
                exit 8
            fi
            echo ""
            wait_until_pods_ready 600 30 $nginx_ns 1
            oc project $install_namespace
        else
            # K8S
            nginx_k8s_yaml="nginx_k8s.yaml"
            cat > ${nginx_k8s_yaml} << __EOF__
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: ${nginx_ns}
  labels:
     app: nginx-stable
spec:
  selector:
    matchLabels:
      app: nginx-stable
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx-stable
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
__EOF__
            kubectl create ns $nginx_ns
            kubectl apply -f $nginx_k8s_yaml
            if [ "$?" != "0" ]; then
                echo -e "\n$(tput setaf 1)Error! create nginx app failed.$(tput sgr 0)"
                leave_prog
                exit 8
            fi
            echo ""
            wait_until_pods_ready 600 30 $nginx_ns 1
        fi
    fi
    echo "Done."
}

add_alamedascaler_for_nginx()
{
    echo -e "\n$(tput setaf 6)Adding nginx alamedascaler ...$(tput sgr 0)"
    nginx_alamedascaler_file="nginx_alamedascaler_file"
    kubectl get alamedascaler -n ${nginx_ns} 2>/dev/null|grep -q "nginx-alamedascaler"
    if [ "$?" != "0" ]; then
        cat > ${nginx_alamedascaler_file} << __EOF__
apiVersion: autoscaling.containers.ai/v1alpha1
kind: AlamedaScaler
metadata:
    name: nginx-alamedascaler
    namespace: ${nginx_ns}
spec:
    policy: stable
    enableExecution: false
    scalingTool:
        type: vpa
    selector:
        matchLabels:
            app: nginx-stable
__EOF__
        kubectl apply -f ${nginx_alamedascaler_file}
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error! Add alamedascaler for nginx app failed.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        sleep 10
    fi
    echo "Done"
}

cleanup_influxdb_prediction_related_contents()
{
    echo -e "\n$(tput setaf 6)Cleaning old influxdb prediction/recommendation/planning records ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    for database in `echo "alameda_prediction alameda_recommendation alameda_planning"`
    do
        echo "database=$database"
        measurement_list="`kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database $database -execute "show measurements" 2>&1 |tail -n+4`"
        for measurement in `echo $measurement_list`
        do
            echo "clean up measurement: $measurement"
            kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database $database -execute "drop measurement $measurement"
        done
    done
    echo "Done."
}

cleanup_alamedaai_models()
{
    #/var/lib/alameda/alameda-ai/models/online/workload_prediction
    echo -e "\n$(tput setaf 6)Cleaning old alameda ai model ...$(tput sgr 0)"
    for ai_pod_name in `kubectl get pods -n $install_namespace -o jsonpath='{range .items[*]}{"\n"}{.metadata.name}'|grep alameda-ai-|grep -v dispatcher`
    do
        kubectl exec $ai_pod_name -n $install_namespace -- rm -rf /var/lib/alameda/alameda-ai/models/online/workload_prediction
    done
    echo "Done."
}

cleanup_influxdb_preloader_related_contents()
{
    echo -e "\n$(tput setaf 6)Cleaning old influxdb preloader metrics records ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    
    measurement_list="`kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "show measurements" 2>&1 |tail -n+4`"
    echo "db=alameda_metric"
    for measurement in `echo $measurement_list`
    do
        if [ "$measurement" = "grafana_config" ]; then
            continue
        fi
        echo "clean up measurement: $measurement"
        kubectl exec $influxdb_pod_name -n $install_namespace -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_metric -execute "drop measurement $measurement"
    done
    
    echo "Done."
}

check_prediction_status()
{
    echo -e "\n$(tput setaf 6)Checking prediction status of monitored objects ...$(tput sgr 0)"
    influxdb_pod_name="`kubectl get pods -n $install_namespace |grep "alameda-influxdb-"|awk '{print $1}'|head -1`"
    measurements_list="`oc exec alameda-influxdb-54949c7c-jp4lk -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_cluster_status -execute "show measurements"|tail -n+4`"
    for measurement in `echo $measurements_list`
    do
        record_number="`oc exec $influxdb_pod_name -- influx -ssl -unsafeSsl -precision rfc3339 -username admin -password adminpass -database alameda_cluster_status -execute "select count(*) from $measurement"|tail -1|awk '{print $NF}'`"
        echo "$measurement = $xx"
        case $future_mode_length in
                ''|*[!0-9]*) echo -e "\n$(tput setaf 1)future mode length (hour) needs to be integer.$(tput sgr 0)" && show_usage ;;
                *) ;;
        esac

        re='^[0-9]+$'
        if ! [[ $xx =~ $re ]] ; then
        echo "error: Not a number" >&2; exit 1
        else
            yy=$(($yy + $xx))
        fi
    done
}

enable_preloader_in_alamedaservice()
{
    get_current_preloader_name
    if [ "$current_preloader_pod_name" != "" ]; then
        echo -e "\n$(tput setaf 6)Skip preloader installation due to preloader pod exist.$(tput sgr 0)"
        echo -e "Delete preloader pod to renew the pod state..."
        kubectl delete pod -n $install_namespace $current_preloader_pod_name
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in deleting preloader pod.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    else
        echo -e "\n$(tput setaf 6)Enable preloader in alamedaservice...$(tput sgr 0)"
        alamedaservice_name="`kubectl get alamedaservice -n $install_namespace -o jsonpath='{range .items[*]}{.metadata.name}'`"
        if [ "$alamedaservice_name" = "" ]; then
            echo -e "\n$(tput setaf 1)Error! Failed to get alamedaservice name.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        kubectl patch alamedaservice $alamedaservice_name -n $install_namespace --type merge --patch '{"spec":{"enablePreloader": true}}'
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in patching alamedaservice $alamedaservice_name.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    fi
    # Check if preloader is ready
    echo ""
    wait_until_pods_ready 600 30 $install_namespace 5
    get_current_preloader_name
    if [ "$current_preloader_pod_name" = "" ]; then
        echo -e "\n$(tput setaf 1)ERROR! Can't find installed preloader pod.$(tput sgr 0)"
        leave_prog
        exit 8
    fi
}

disable_preloader_in_alamedaservice()
{
    echo -e "\n$(tput setaf 6)Disable preloader in alamedaservice...$(tput sgr 0)"
    get_current_preloader_name
    if [ "$current_preloader_pod_name" != "" ]; then
        alamedaservice_name="`kubectl get alamedaservice -n $install_namespace -o jsonpath='{range .items[*]}{.metadata.name}'`"
        if [ "$alamedaservice_name" = "" ]; then
            echo -e "\n$(tput setaf 1)Error! Failed to get alamedaservice name.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
        kubectl patch alamedaservice $alamedaservice_name -n $install_namespace  --type merge --patch '{"spec":{"enablePreloader": false}}'
        if [ "$?" != "0" ]; then
            echo -e "\n$(tput setaf 1)Error in patching alamedaservice $alamedaservice_name.$(tput sgr 0)"
            leave_prog
            exit 8
        fi

        # Check if preloader is removed and other pods are ready
        echo ""
        wait_until_pods_ready 600 30 $install_namespace 5
        get_current_preloader_name
        if [ "$current_preloader_pod_name" != "" ]; then
            echo -e "\n$(tput setaf 1)ERROR! Can't stop preloader pod.$(tput sgr 0)"
            leave_prog
            exit 8
        fi
    fi
    echo "Done"
}

clean_environment_operations()
{
    cleanup_influxdb_preloader_related_contents
    cleanup_influxdb_prediction_related_contents
    cleanup_alamedaai_models
}   

if [ "$#" -eq "0" ]; then
    show_usage
    exit
fi

while getopts "f:ecpvrdh" o; do
    case "${o}" in
        p)
            prepare_environment="y"
            ;;
        c)
            clean_environment="y"
            ;;
        e)
            enable_preloader="y"
            ;;
        r)
            run_preloader="y"
            ;;
        f)
            future_mode_enabled="y"
            f_arg=${OPTARG}
            ;;
        d)
            disable_preloader="y"
            ;;
        v)
            revert_environment="y"
            ;;
        h)
            show_usage
            exit
            ;;
        *)
            echo "Warning! wrong paramter, ignore it."
            ;;
    esac
done

if [ "$future_mode_enabled" = "y" ]; then
    future_mode_length=$f_arg
    case $future_mode_length in
        ''|*[!0-9]*) echo -e "\n$(tput setaf 1)future mode length (hour) needs to be integer.$(tput sgr 0)" && show_usage ;;
        *) ;;
    esac
fi

kubectl version|grep -q "^Server"
if [ "$?" != "0" ];then
    echo -e "\nPlease login to kubernetes first."
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

file_folder="/tmp/preloader"
nginx_ns="nginx-preloader-sample"

rm -rf $file_folder
mkdir -p $file_folder
current_location=`pwd`
cd $file_folder

if [ "$prepare_environment" = "y" ]; then
    delete_all_alamedascaler
    new_nginx_example
    patch_datahub_for_preloader
    patch_grafana_for_preloader
    check_influxdb_retention
    add_alamedascaler_for_nginx
fi

if [ "$clean_environment" = "y" ]; then
    clean_environment_operations
fi

if [ "$enable_preloader" = "y" ]; then
    enable_preloader_in_alamedaservice
fi

if [ "$run_preloader" = "y" ]; then
    # Enable preloader will reset the datahub env. Need to set env for datahub again.
    patch_datahub_for_preloader
    # clean up again
    clean_environment_operations

    run_preloader_command
    verify_metrics_exist
    reschedule_dispatcher
    #check_prediction_status
fi

if [ "$future_mode_enabled" = "y" ]; then
    # Enable preloader will reset the datahub env. Need to set env for datahub again.
    patch_datahub_for_preloader
    run_futuremode_preloader
    verify_metrics_exist
fi

if [ "$disable_preloader" = "y" ]; then
    disable_preloader_in_alamedaservice
fi

if [ "$revert_environment" = "y" ]; then
    delete_all_alamedascaler
    delete_nginx_example
    patch_datahub_back_to_normal
    patch_grafana_back_to_normal
    clean_environment_operations
fi

leave_prog
exit 0
