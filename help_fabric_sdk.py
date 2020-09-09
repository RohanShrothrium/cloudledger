import asyncio
from hfc.fabric import Client
loop = asyncio.get_event_loop()
cli = Client(net_profile="/home/rohan/Documents/cloudledger/network.json")
org1_admin = cli.get_user(org_name='org1.example.com', name='Admin')
peer_target=cli._peers['peer0.org1.example.com']
requestor=org1_admin


loop.run_until_complete(cli.init_with_discovery(org1_admin, peer_target, 'mychannel'))

args=['170020031', 'rohan', 'rohan', 'rohan@g.com', '100']


response = loop.run_until_complete(cli.chaincode_invoke(
                    requestor=org1_admin,
                    channel_name='mychannel',
                    peers=peers,
                    args=args,
                    cc_name=cc_name,
                    transient_map=None,
                    wait_for_event=True
                    ))


responses = loop.run_until_complete(cli.chaincode_install(
               requestor=org1_admin,
               peers=['peer0.org1.example.com',
                      'peer1.org1.example.com'],
               cc_path='github.com/example_cc',
               cc_name='mycc',
               cc_version='v1.0'
               ))


response = loop.run_until_complete(cli.chaincode_instantiate(
               requestor=org1_admin,
               channel_name='mychannel',
               peers=['peer0.org1.example.com'],
               args=args,
               cc_name='mycc',
               cc_version='1.0',
               collections_config=None,
               wait_for_event=True # optional, for being sure chaincode is instantiated
               ))


response = loop.run_until_complete(cli.chaincode_query(
               requestor=org1_admin,
               channel_name='mychannel',
               peers=['peer0.org1.example.com'],
               fcn='Query',
               args=['170020031'],
               cc_name='mycc'
               ))