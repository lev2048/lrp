import React, { useState } from 'react';
import { BsPlusCircleDotted } from "react-icons/bs";
import { BiPlanet, BiTrash, BiTerminal } from "react-icons/bi";
import { ProxyModal } from "../ProxyModal";
import { ProxyDelModal } from "../ProxyDelModal";
import "./index.css";

interface ProxyInfo {
    id: string
    cid: string
    info: string
    mark: string
    status: number
    isTemp: boolean
    connNum: number
}

interface Client {
    id: string;
    mark: string;
    online: boolean;
    proxyInfos: ProxyInfo[];
}

interface IProps {
    data: Client[]
}

const ProxyList: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    const [showAddModal, setAddModal] = useState(false);
    const [showDelModal, setDelModal] = useState(false);
    const [delProxyId, setDelProxyId] = useState({ pid: "", cid: "" });

    let proxyList: ProxyInfo[] = [];
    props.data.forEach((c: Client) => {
        if (c.online) {
            let data: ProxyInfo[] = [];
            c.proxyInfos.forEach(v => {
                v.cid = c.id;
                data.push(v);
            })
            proxyList = proxyList.concat(data);
        }
    });

    let content: JSX.Element[] = proxyList.map((v, k) => (
        <div className="proxyListItem" key={k}>
            <div className="proxyIcon">{v.isTemp ? <BiTerminal /> : <BiPlanet />}</div>
            <div style={{ width: '21%' }}>{v.mark !== "" ? v.mark : v.id}</div>
            <div style={{ width: '30%' }}>{v.info}</div>
            <div style={{ width: '10%' }}>{v.connNum} conn</div>
            <div className="proxyStatus">{v.status === 1 ? "Running" : "Warning"}</div>
            <div className="proxyDel" onClick={() => { setDelModal(true); setDelProxyId({ pid: v.id, cid: v.cid }) }}><BiTrash /></div>
        </div>
    ));

    return (
        <div className="proxyList">
            <div className="proxyListHeader">
                <div className="InfoTitle">ProxyList</div>
                <div className="proxyAdd" onClick={() => setAddModal(!showAddModal)}><BsPlusCircleDotted /></div>
            </div>
            <div className="proxyListContent">
                {content.length !== 0 ? content : (<div className="emptyContent"><img src="/images/empty.png" alt="" /></div>)}
            </div>
            {showAddModal && <ProxyModal setModal={setAddModal} />}
            {showDelModal && <ProxyDelModal setModal={setDelModal} pid={delProxyId.pid} cid={delProxyId.cid} />}
        </div>
    );
}

export default ProxyList;