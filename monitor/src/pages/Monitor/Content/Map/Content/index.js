import './index.css';

export default function Content({map}){
    return (
        <div className="monitor-map-content">
            <img src={map?.url} alt='' />
        </div>
    );
}