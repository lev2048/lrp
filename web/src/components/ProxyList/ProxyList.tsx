import React from 'react';
import { BsPlusCircleDotted } from "react-icons/bs";
import { BiPlanet, BiTrash, BiTerminal } from "react-icons/bi";
import "./index.css";

interface ProxyInfo {
    cid: string
    pid: string
    info: string
    mark: string
    status: number
    isTemp: boolean
    connNum: number
}

interface IProps {
    data: ProxyInfo[]
}

const ProxyList: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    let content: JSX.Element[] = props.data.map((v, k) => (
        <div className="proxyListItem" key={k}>
            <div className="proxyIcon">{v.isTemp ? <BiTerminal /> : <BiPlanet />}</div>
            <div style={{ width: '21%' }}>{v.mark !== "" ? v.mark : v.pid}</div>
            <div style={{ width: '30%' }}>{v.info}</div>
            <div style={{ width: '10%' }}>{v.connNum} conn</div>
            <div className="proxyStatus">{v.status === 1 ? "Running" : "Warning"}</div>
            <div className="proxyDel"><BiTrash /></div>
        </div>
    ))
    return (
        <div className="proxyList">
            <div className="proxyListHeader">
                <div className="InfoTitle">ProxyList</div>
                <div className="proxyAdd"><BsPlusCircleDotted /></div>
            </div>
            <div className="proxyListContent">
                {content && (<div className="emptyContent"><img src="/images/empty.png" alt="" /></div>)}
            </div>
        </div>
    );
}

export default ProxyList;