'use strict';
var log4js = require('log4js');
var logger = log4js.getLogger('SampleWebApp');
var express = require('express');
var session = require('express-session');
var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');
var http = require('http');
var util = require('util');
var helper = require('./app/helper.js');
var joinChannel = require('./app/join-channel.js');
var createChannel = require('./app/create-channel.js');
var installChainCode = require('./app/install-chaincode.js');
require('./config.js');
var hfc = require('fabric-client');
async function printTest(){
    var res = await helper.getRegisteredUser('user1','bob',true).then(function(res) {return res});
    console.log(res);
    
}

async function testCreateChannel(){
    // var res = await createChannel.createChannel('transfers','../transfers.tx','user1','alice').then(function(res){return res})
    var res = await createChannel.createChannel('fredrick-alice','../fredrick-alice.tx','user1','alice').then(function(res){return res})
}

async function testJoinChannel(){
    var res = await joinChannel.joinChannel("fredrick-alice",["peer0.alice.coderschool.vn"],"user1","alice")
}

async function testJoinChannel(){
    var res = await joinChannel.joinChannel("fredrick-alice",["peer0.alice.coderschool.vn"],"user1","alice")
}

/**
 * Get peer host:port by composite id.
 * Composite ID is a combination of orgID and peerID, split by '/'. For example: 'org1/peer2'
 * @param {string} orgPeerID
 * @returns {string}
 */
function getPeerHostByCompositeID(orgPeerID){
    var parts = orgPeerID.split('/');
    var peer = networkConfig[parts[0]][parts[1]] || {};
    return tools.getHost(peer.requests);
}
/**
 * Get peer host:port by composite id.
 * Composite ID is a combination of orgID and peerID, split by '/'. For example: 'org1/peer2'
 * @param {string} orgPeerID
 * @returns {{peer:string, org:string}}
 */
function getPeerInfoByCompositeID(orgPeerID){
    var parts = orgPeerID.split('/');
    return parts ? {peer: parts[1], org: parts[0]} : null;
}

async function testInstallChaincode(){
    logger.debug('==================== INSTALL CHAINCODE ==================');

    var chaincodeName = "SalmonRecord";
    var chaincodePath = "github.com/example_cc/go";
    var chaincodeVersion = "v0";
    var peersId = ["peer0.alice.coderschool.vn","peer0.bob.coderschool.vn"];
    var peers   = peersId.map(getPeerHostByCompositeID);

    logger.debug('peers : ' + peers); // target peers list
    logger.debug('chaincodeName : ' + chaincodeName);
    logger.debug('chaincodePath  : ' + chaincodePath);
    logger.debug('chaincodeVersion  : ' + chaincodeVersion);
    if (!peers || peers.length === 0) {
        res.error(getErrorMessage('\'peers\''));
        return;
    }
    if (!chaincodeName) {
        res.error(getErrorMessage('\'chaincodeName\''));
        return;
    }
    if (!chaincodePath) {
        res.error(getErrorMessage('\'chaincodePath\''));
        return;
    }
    if (!chaincodeVersion) {
        res.error(getErrorMessage('\'chaincodeVersion\''));
        return;
    }

    res.promise(
      install.installChaincode(peers, chaincodeName, chaincodePath, chaincodeVersion, USERNAME, ORG)
    );
    var res = await installChainCode.installChaincode("fredrick-alice",["peer0.alice.coderschool.vn"],"user1","alice")
}
printTest();

// testCreateChannel();
// testJoinChannel();