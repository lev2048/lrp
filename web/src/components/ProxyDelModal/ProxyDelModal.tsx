import React from 'react';
import "./index.css";
import axios from 'axios';
import { AiOutlineDelete } from "react-icons/ai";

interface IProps {
    cid: string;
    pid: string;
    setModal: (value: boolean) => void;
}

const ProxyDelModal: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    let postDelReq = () => {
        let data = new FormData();
        data.append('cid', props.cid);
        data.append('pid', props.pid);
        axios.post("http://10.0.0.216/v1/proxy/del", data).then(res => {
            console.log(res.data);
            props.setModal(false);
        });
    }
    return (
        <div className="modalMask">
            <div className="animate__animated animate__zoomIn animate__faster proxyDelModal">
                <div className="proxyDelIcon"><AiOutlineDelete /></div>
                <div className="proxyDelTitle">Delete Proxy</div>
                <div className="proxyDelInfo">This action can't be undone.</div>
                <div className="proxyDelButton">
                    <div className="proxyDelCancel" onClick={() => props.setModal(false)}>Cancel</div>
                    <div className="proxyDelAction" onClick={() => postDelReq()}>Delete</div>
                </div>
            </div>
        </div>
    );
}

export default ProxyDelModal;