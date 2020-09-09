#!/bin/bash
function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp { 
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$5#" \
        -e "s#\${CAPEM}#$6#" \
        ccp-template.json 
}
for((i=1; i<=$1; i++))
do
    PORT_NOS="$(docker ps -a | grep hyperledger/fabric-peer | grep org$i | awk '{print $(NF-2)}')"
    for((j=1; j<=$2; j++))
    do
        eval P$(($j-1))PORT=$(for p in $PORT_NOS ; do echo $p | cut -d '-' -f 1 | cut -d':' -f 2;  done | awk 'NR=='$j'{print $1}';)
    done 
    ORG=$i
    CAPORT=$(docker ps -a | grep hyperledger/fabric-ca | egrep -h "Org$i|org$i|ORG$i"  | awk '{print $(NF-1)}' | cut -d '>' -f 2 | cut -d '/' -f 1)
    PEERPEM=crypto-config/peerOrganizations/org$i.example.com/tlsca/tlsca.org$i.example.com-cert.pem
    CAPEM=crypto-config/peerOrganizations/org$i.example.com/ca/ca.org$i.example.com-cert.pem
    echo $ORG $P0PORT $CAPORT 
    echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM)" > connection-org$i.json

done
