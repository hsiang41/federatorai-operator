#!/usr/bin/env bash

#################################################################################################################
#
#   This script is created for assign node label for cost analysis function (lab usage)
#
#################################################################################################################

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
        echo -e "\n$(tput setaf 10)Error! OpenShift version less than 3.$openshift_required_minor_version is not supported by Federator.ai$(tput sgr 0)"
        exit 5
    fi

    # oc major version = 3
    openshift_minor_version=`oc version 2>/dev/null|grep "oc v"|cut -d '.' -f2`
    # k8s version = 1.x
    k8s_version=`kubectl version 2>/dev/null|grep Server|grep -o "Minor:\"[0-9]*\""|cut -d '"' -f2`

    if [ "$openshift_minor_version" != "" ] && [ "$openshift_minor_version" -lt "$openshift_required_minor_version" ]; then
        echo -e "\n$(tput setaf 10)Error! OpenShift version less than 3.$openshift_required_minor_version is not supported by Federator.ai$(tput sgr 0)"
        exit 5
    elif [ "$openshift_minor_version" = "" ] && [ "$k8s_version" != "" ] && [ "$k8s_version" -lt "$k8s_required_version" ]; then
        echo -e "\n$(tput setaf 10)Error! Kubernetes version less than 1.$k8s_required_version is not supported by Federator.ai$(tput sgr 0)"
        exit 6
    elif [ "$openshift_minor_version" = "" ] && [ "$k8s_version" = "" ]; then
        echo -e "\n$(tput setaf 10)Error! Can't get Kubernetes or OpenShift version$(tput sgr 0)"
        exit 5
    fi
}

assign_label_to_each_nodes()
{

    while read node_name node_status node_role others
    do
        echo -e "\n$(tput setaf 6)Starting label $node_name : $node_role node ...$(tput sgr 0)"
        if [ "$node_name" = "" ] || [ "$node_role" = "" ]; then
           continue
        fi

        if [[ $node_role == *"master"* ]]; then
            kubectl label --overwrite node $node_name node-role.kubernetes.io/master=""
        fi

        kubectl label node $node_name beta.kubernetes.io/instance-type=c5.xlarge
        kubectl label node $node_name failure-domain.beta.kubernetes.io/region=us-west-2
        kubectl label node $node_name failure-domain.beta.kubernetes.io/zone=us-west-2a
        kubectl patch nodes $node_name --type merge --patch '{"spec":{"providerID": "aws:///us-west-2a/i-0b186a4958f5f8576"}}'
        echo "Done"
    done <<< "$(kubectl get nodes|grep -v STATUS)"
}

kubectl version|grep -q "^Server"
if [ "$?" != "0" ];then
    echo -e "\nPlease login to kubernetes first."
    exit
fi

echo "Checking environment version..."
check_version
echo "...Passed"

assign_label_to_each_nodes

exit 0
