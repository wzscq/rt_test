import {useSelector} from 'react-redux';
import Indicator from './Indicator';

import './index.css';

export default function Content({map}){
    const indicator=useSelector(state=>state.data.indicator);
    const points=useSelector(state=>state.data.points);
    const currentPos=useSelector(state=>state.data.currentPos);
    
    const pointsControl=points.map((dataItem,index)=>{
        const {x,y,rgb,value}=dataItem;
        const isCurPoint=(index+1===currentPos)?true:false;
        if(isCurPoint===true){
            return (
                <>    
                    <div className='map-point-label' style={{left:x,top:y-25}}>{value}</div>
                    <div className='map-point' style={{left:x,top:y,backgroundColor:rgb}}></div>
                </>
            );
        }

        return (
            <div className='map-point' style={{left:x,top:y,backgroundColor:rgb}}></div>
        );
    });

    return (
        <div className="monitor-map-content">
            <Indicator indicator={indicator}/>
            <img src={map?.url} alt='' />
            {pointsControl}
        </div>
    );
}