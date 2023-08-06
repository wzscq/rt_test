import {useSelector} from 'react-redux';
import { SplitPane } from "react-collapse-pane";
import PropertyGrid from './PropertyGrid';
import Map from './Map';

import './index.css';

export default function Content({sendMessageToParent}){
  const data=useSelector(state=>state.data.data.length>0?state.data.data[state.data.data.length-1]:undefined);
  return (
    <div className='monitor-content'>
      <SplitPane dir='ltr'initialSizes={[65,35]} split="vertical" collapse={false}>
        <div className='monitor-content-left'>
        <SplitPane dir='rtl'initialSizes={[60,40]} split="horizontal" collapse={false}>
          <Map sendMessageToParent={sendMessageToParent}/>    
          <div></div>
        </SplitPane>
        </div>
        <div className='monitor-content-right'>
          <SplitPane dir='rtl'initialSizes={[20,20,20,20,10]} split="horizontal" collapse={false}>
            <PropertyGrid obj={data?.radio?.measures_common} title="common measures"/>
            <PropertyGrid obj={data?.radio?.measures_lte} title="lte measures"/>
            <PropertyGrid obj={data?.radio?.measures_nr} title="nr measures"/>
            <PropertyGrid obj={data?.robot_info} title="robot info"/>
            <PropertyGrid obj={data?.case_progress} title="case progress"/>
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  )
}