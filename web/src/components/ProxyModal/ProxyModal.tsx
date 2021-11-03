import React, { useState } from 'react';
import axios from 'axios';
import "./index.css";
import 'animate.css';

interface IProps {
    setModal: (value: boolean) => void
}

const ProxyModal: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    const [req, setReqData] = useState({
        mark: "",
        clientId: "",
        destAddr: "",
        listenPort: "",
    });

    let postAddReq = () => {
        console.log(req);
        let data = new FormData();
        data.append('mark', req.mark);
        data.append('clientId', req.clientId);
        data.append('destAddr', req.destAddr);
        data.append('listenPort', req.listenPort);
        axios.post("http://10.0.0.216/v1/proxy/add", data).then(res => {
            console.log(res.data);
            props.setModal(false);
        });
    }

    return (
        <div className="modalMask">
            <div className="animate__animated animate__zoomIn animate__faster proxyModal">
                <div className="modalTitle">Create Proxy</div>
                <div className="modalInput">
                    <div className="modalInputTitle">Name</div>
                    <input type="text" onChange={e => { req.mark = e.target.value; setReqData(req) }} placeholder="Like Home" />
                </div>
                <div className="modalInput">
                    <div className="modalInputTitle">ClinetID</div>
                    <input type="text" onChange={e => { req.clientId = e.target.value; setReqData(req) }} placeholder="Like c5tuk9b765s2jfnbtfog" />
                </div>
                <div className="modalInput">
                    <div className="modalInputTitle">DestAddr</div>
                    <input type="text" onChange={e => { req.destAddr = e.target.value; setReqData(req) }} placeholder="Like 10.0.0.1:80" />
                </div>
                <div className="modalInput">
                    <div className="modalInputTitle">ListenPort</div>
                    <input type="number" onChange={e => { req.listenPort = e.target.value; setReqData(req) }} placeholder="Like 3389" />
                </div>
                <div className="modalButton">
                    <div className="modalCancel" onClick={() => props.setModal(false)}>Cancel</div>
                    <div className="modalOK" onClick={() => postAddReq()}>Create</div>
                </div>
            </div>
        </div>
    );
}

export default ProxyModal;