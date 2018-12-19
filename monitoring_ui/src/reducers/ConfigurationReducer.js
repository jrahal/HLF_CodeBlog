/*****************************************************************************
Copyright (c) 2016 IBM Corporation and other Contributors.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.


Contributors:

Alex Nguyen - Initial Contribution
*****************************************************************************/
import {
  CONFIG_UPDATE_URL_REST_ROOT, SET_CONFIG_IOT_CONNECTION, SET_CHAIN_HEIGHT_POLLING_INTERVAL_ID, SET_OBC_CONFIGURATION,
  SET_CONFIG_DIALOG_DISPLAY
} from '../actions/ConfigurationActions.js'

//set default configuration. To keep things simple, we don't have any UI elements
//to configure these properties, so they must be populated ahead of time.
export const configuration = (state={
  urlRestRoot: "https://ibmblockchain-starter.ng.bluemix.net/api/v1",
   //urlRestRoot: "http://169.44.63.199:37687",
  //chaincodeId: "abf072028033b86aa9d61127c9ac9f0f407a24ce5b464f3afb6e8474169df95e1c1e40d31051553430eca22d10fe8b7083a518c77b80c94d679dba4c6858a90b",
  chaincodeId: "simple_contract",
  //the ID of the chain height polling. This will be populated at runtime.
  chainHeightPollingIntervalId: -1,
  //the chain height polling interval, in milliseconds
  chainHeightPollingInterval: 2000,
  secureContext: "user_context",
  //this is a UI specific toggle. This governs whether or not the configuration modal is showing or not.
  showDialog: false,
  //this is the number of blocks to show on the page at a time.
  blocksPerPage: 10,
  key: "org1",
  secret: "secret",
  networkId: "networkId",
  channel: "defaultchannel",
  iot_org: "",
  iot_api_key: "",
  iot_auth_token: ""
}, action) => {
  /**
    These state configuration actions are implemented, but we don't use them in the UI. Feel free to wire them to a UI element if needed.
  **/
  switch (action.type){
    case CONFIG_UPDATE_URL_REST_ROOT:
      return Object.assign({}, state, {
        urlRestRoot: action.url
      })
    case SET_CHAIN_HEIGHT_POLLING_INTERVAL_ID: {
      return Object.assign({}, state, {
        chainHeightPollingIntervalId: action.intervalId
      })
    }
    /**
      This action is called from the submit button on the form. It is used to transfer over values from the obcConfiguration model to the configuration used by the ui.
      This allows us to make changes to the form without affecting the queries going on in the background.
    **/
    case SET_OBC_CONFIGURATION:
      var params = {
        //set the appropriate properties
        urlRestRoot: action.obcConfigObj.urlRestRoot,
        chaincodeId: action.obcConfigObj.chaincodeId,
        secureContext: action.obcConfigObj.secureContext,
        blocksPerPage: action.obcConfigObj.blocksPerPage,
        key: action.obcConfigObj.key,
        secret: action.obcConfigObj.secret,
        networkId: action.obcConfigObj.networkId,
        channel: action.obcConfigObj.channel,
        iotOrg: action.obcConfigObj.iotOrg,
        iotAuthToken: action.obcConfigObj.iotAuthToken,
        iotApiKey: action.obcConfigObj.iotApiKey
      }
      let config = {
          method: 'POST',
          headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json'
          },
          body: JSON.stringify(params)
      }
      console.log("updating client")
      fetch('http://localhost:3001/init_client', config).then( (res) => { fetch('http://localhost:3001/getchaincodes', config) } )
      return Object.assign({}, state, params)

    /**
      Strictly a UI control. This determines whether or not the configuration ui dialog should display or not.
    **/
    case SET_CONFIG_DIALOG_DISPLAY:
      return Object.assign({}, state, {
        showDialog: action.showDialog
      })

    default:
      return state;
  }
}
