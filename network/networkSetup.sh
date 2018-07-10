#!/bin/bash

export PATH=${PWD}/tools:${PWD}:$PATH
export FABROC_CFG_PATH=${PWD}
export COMPOSE_PROJECT_NAME="tictactoe"
export IMAGE_TAG="latest"

function printHelp () {
	echo "${PWD}"
	# echo $PATH
	echo "Usage: "
	echo "This is still under construction"
}

function clearContainers () {
	CONTAINER_IDS=$(docker ps -aq)
	if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
		echo "---- No containers availalble for deletion ----"
	else 
		docker rm -f $CONTAINER_IDS
	fi
}

function removeUnwantedImages () {
	DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]" | awk '{print $3}')
	if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
		echo "---- No images availalble for deletion ----"
	else
		docker rmi -f $DOCKER_IMAGE_IDS
	fi
}

BLACKLISTED_VERSIONS="^1\.0\. ^1\.1\.0-preview ^1\.1\.0-alpha"

function checkPrereqs () {
    echo ${PATH}
	LOCAL_VERSION=$(configtxlator version | sed -ne 's/ Version: //p')
	DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-tools:$IMAGETAG peer version | sed -ne 's/ Version: //p' | head -1)

	echo "LOCAL_VERSION=$LOCAL_VERSION"
	echo "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

	if [[ $LOCAL_VERSION != $DOCKER_IMAGE_VERSION ]]; then
		echo "========================== Warning ========================="
		echo "  Local fabric binaries and docker images are out of sync.  "
		echo "  This may cause problems"
		echo "============================================================"
	fi

	for UNSUPPORTED_VERSION in $BLACKLISTED_VERSIONS ; do
		echo "$LOCAL_VERSION" | grep -q $UNSUPPORTED_VERSION
		if [[ $? -eq 0 ]]; then
			echo "ERROR! Local fabric binary version of $LOCAL_VERSION does not match this newer version of byfn and is unsupported. Move to a later version or checkout earlier version of fabric-samples"
			exit 1
		fi

		echo "$DOCKER_IMAGE_VERSION" | grep -q $UNSUPPORTED_VERSION
		if [ $? -eq 0 ]; then
			echo "ERROR! Fabric docker image version $DOCKER_IMAGE_VERSION does not match this newer version of byfn and is unsupported. Move to later version of checkout earlier version of fabric-samples"
			exit 1
		fi
	done
}

function askProceed () {
	read -p "Continue? [y/n] " ans
	case "$ans" in
		y|Y|"" )
			echo "proceeding..."
		;;
		n|N )
			echo "exiting..."
			exit 1
		;;
		* )
			echo "invalid response"
			askProceed
		;;
	esac
}

function replacePrivateKey () {
	ARCH=`uname -s | grep Darwin`
	if [[ $ARCH == "Darwin" ]]; then
		OPTS="-it"
	else
			OPTS="-i"
	fi

	cp docker-compose-e2e-template.yaml docker-compose-e2e.yaml

	CURRENT_DIR=$PWD
	cd crypto-config/peerOrganizations/org1.tictactoe.com/ca/
	PRIV_KEY=$(ls *_sk)
	cd "$CURRENT_DIR"
	sed $OPTS "s/ORG1CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

	cd crypto-config/peerOrganizations/org2.tictactoe.com/ca/
	PRIV_KEY=$(ls *_sk)
	cd "$CURRENT_DIR"
	sed $OPTS "s/ORG2CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-e2e.yaml

	if [[ $ARCH == "Darwin" ]]; then
		rm docker-compose-e2e.yaml
	fi
}

function networkUp () {
	checkPrereqs

	if [ ! -d "crypto-config" ]; then
		generateCerts
		replacePrivateKey
		generateChannelArtifacts
	fi

	if [ "${IF_COUCHDB}" == "couchdb" ]; then
		IMAGE_TAG=$IMAGETAG docker-compose -f $COMPOSE_FILE -f $COMPOSE_FILE_COUCH up -d 2>&1
	else
		IMAGE_TAG=$IMAGETAG docker-compose -f $COMPOSE_FILE up -d 2>&1
	fi

	if [ $? -ne 0 ]; then
		echo "ERROR !!! unable to start network"
		exit 1
	fi

	#docker exec cli scripts/script.sh $CHANNEL_NAME $CLI_DELAY $LANGUAGE $CLI_TIMEOUT
	#if [ $? -ne 0 ]; then
	#	echo "ERROR !!! Test failed"
	#	exit 1
	#fi
}

function networkDown () {
	docker-compose -f $COMPOSE_FILE -f $COMPOSE_FILE_COUCH down --volumes
	docker-compose -f $COMPOSE_FILE down --volumes

	if [ "$MODE" != "restart" ]; then
		docker run -v $PWD:/tmp/first-network -rm hyperledger/fabric-tools:$IMAGETAG rm -Rf /tmp/first-network/ledger-backup
		clearContainers
		removeUnwantedImages
#		 rm -rf channel-artifacts/*.block channel-artifacts/*.tx crypto-config
#		 rm -f docker-compose-e2e.yaml
	fi
}

function generateCerts () {	
	which cryptogen
	if [ "$?" -ne 0 ]; then
		echo "cryptogen tool not found"
		exit 1
	fi
	echo
	echo "############################################################################"
	echo "####### generate certificates using cryptogen and crypto-config.yaml #######"
	echo "############################################################################"

	if [ -d "crypto-config" ]; then
		rm -Rf crypto-config
	fi
	set -x
	cryptogen generate --config=./crypto-config.yaml
	res=$?
	set +x
	if [ $res -ne 0 ]; then
		echo "failed to generate certificates"
		exit 1
	fi
	echo
}

function generateChannelArtifacts () {
	which configtxgen
	if [ "$?" -ne 0 ]; then
		echo "configtxgen tool not found. exiting"
		exit 1
	fi

	echo
	echo "############################################################################"
	echo "####### generate channel artifacts using cryptogen and configtx.yaml #######"
	echo "############################################################################"

	echo
	echo "##############################################"
	echo "####### generate orderer genesis block #######"
	echo "##############################################"

	set -x
	configtxgen -profile TwoOrgOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
	res=$?
	set +x
	if [ $res -ne 0 ]; then
		echo "failed to generate genesis block"
		exit 1
	fi

# channel will be created by the application

#	echo
#	echo "#######################################################################"
#	echo "####### generate channel configuration transaction 'channel.tx' #######"
#	echo "#######################################################################"
#
#	set -x
#	configtxgen -profile TicTacToeChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
#	res=$?
#	set +x
#	if [ $res -ne 0 ]; then
#		echo "failed to create channel transaction block"
#		exit 1
#	fi
#
#	echo
#	echo "########################################################"
#	echo "####### generate anchor peer update for sivatech #######"
#	echo "########################################################"
#
#	set -x
#	configtxgen -profile TicTacToeChannel -outputAnchorPeersUpdate ./channel-artifacts/TicTacToeMSPAnchors.tx -channelID $CHANNEL_NAME -asOrg TicTacToeMSP
#	res=$?
#	set +x
#	if [ $res -ne 0 ]; then
#		echo "failed to create achor peer update for sivatech"
#		exit 1
#	fi

	echo
}

OS_ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')

CLI_TIMEOUT=10

CLI_DELAY=3

CHANNEL_NAME="t3-sivatech"

COMPOSE_FILE=docker-compose-e2e.yaml

COMPOSE_FILE_COUCH=docker-compose-couch.yaml

LANGUAGE=golang

IMAGETAG="latest"

# Parse command line args
if [ "$1" = "-m" ]; then
	shift
fi

MODE=$1;shift

# Determine operation mode
if [ "$MODE" == "up" ]; then
	EXPMODE="starting"
elif [ "$MODE" == "down" ]; then
	EXPMODE="teardown"
elif [ "$MODE" == "generatecert" ]; then
	EXPMODE="generatecert"
elif [ "$MODE" == "generateartifacts" ]; then
	EXPMODE="generateartifacts"
elif [[ $MODE == "checkPrereq" ]]; then
	EXPMODE="checkPrereq"
elif [[ $MODE == "clearContainers" ]]; then
	EXPMODE="clearContainers"
elif [[ $MODE == "replacePrivateKey" ]]; then
	EXPMODE="replacePrivateKey"
else
	printHelp
	exit 1
fi

while getopts "h?m:c:t:d:f:s:l:i:" opt
do
	case "$opt" in
		h|\?)
			printHelp
			exit 0
		;;
		c) CHANNEL_NAME=$OPTARG
		;;
		i) IMAGETAG=`uname -m`"-"$OPTARG
		;;
	esac
done

askProceed

if [ "${MODE}" == "generatecert" ]; then
	generateCerts
elif [ "${MODE}" == "generateartifacts" ]; then
	generateChannelArtifacts
elif [[ $MODE == "checkPrereq" ]]; then
	checkPrereqs
elif [[ $MODE == "clearContainers" ]]; then
	clearContainers
elif [[ $MODE == "replacePrivateKey" ]]; then
	replacePrivateKey
elif [[ $MODE == "up" ]]; then
    networkUp
elif [[ $MODE == "down" ]]; then
    networkDown
else
	printHelp
	exit 1
fi
