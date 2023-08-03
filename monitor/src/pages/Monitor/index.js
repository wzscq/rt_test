import {useSelector} from 'react-redux';
import mqtt from 'mqtt';
import { useEffect } from 'react';

var g_MQTTClient=null;

export default function Monitor(){
  const {mqttConf}=useSelector(state=>state.mqtt);

  const connectMqtt=()=>{
    console.log("connectMqtt ... ");
    if(g_MQTTClient!==null){
        g_MQTTClient.end();
        g_MQTTClient=null;
    }

    const server='ws://'+mqttConf.broker+':'+mqttConf.wsPort;
    const options={
        username:mqttConf.user,
        password:mqttConf.password,
    }
    console.log("connect to mqtt server ... "+server+" with options:",options);
    g_MQTTClient  = mqtt.connect(server,options);
    g_MQTTClient.on('connect', () => {
        console.log("connected to mqtt server "+server+".");
        console.log("subscribe topics ...");
        g_MQTTClient.subscribe("realtime_measurement_reporting/178BFBFF00800F82", (err) => {
            if(!err){
                console.log("subscribe topics success.");
                console.log("topic:","realtime_measurement_reporting/178BFBFF00800F82");
                //发送流执行请求
             
            } else {
                console.log("subscribe topics error :"+err.toString());
            }
        });
    });
    g_MQTTClient.on('message', (topic, payload, packet) => {
        console.log("receiconsolleve message topic :"+topic+" content :"+payload.toString());
    });
    g_MQTTClient.on('close', () => {
        console.log("mqtt client is closed.");
    });
  }

  useEffect(()=>{
    connectMqtt();
  })

  return (
    <div>Monitor</div>
  )
}