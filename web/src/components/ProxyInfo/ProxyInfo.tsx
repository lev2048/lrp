import React from 'react';
import "./index.css";

interface IProps {
    connNum: number
    proxyNum: number
    tproxyNum: number
    clientNum: number
}

const ProxyInfo: React.FunctionComponent<IProps> = (props:IProps): JSX.Element => {
    return (
        <div className="proxyInfo">
            <div className="proxyItem">
                <div className="proxyItemTitle">Client</div>
                <div>{props.clientNum}</div>
            </div>
            <div className="splitLine"></div>
            <div className="proxyItem">
                <div className="proxyItemTitle">Proxy</div>
                <div>{props.tproxyNum}/{props.proxyNum}</div>
            </div>
            <div className="splitLine"></div>
            <div className="proxyItem">
                <div className="proxyItemTitle">Conn</div>
                <div>{props.connNum}</div>
            </div>
        </div>
    );
}

export default ProxyInfo;