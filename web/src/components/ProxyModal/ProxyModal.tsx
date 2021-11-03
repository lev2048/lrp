import React from 'react';
import "./index.css";
import 'animate.css';

interface IProps {
    setModal: (value: boolean) => void
}

const ProxyModal: React.FunctionComponent<IProps> = (props: IProps): JSX.Element => {
    return (
        <div className="modalMask">
            <div className="animate__animated animate__zoomIn animate__faster proxyModal">
                <div className="modalTitle">Create Proxy</div>
                <div className="modalInput">
                    <div className="modalInputTitle">Name</div>
                    <input type="text" placeholder="Like Home" />
                </div>
                <div className="modalInput">
                    <div className="modalInputTitle">ClinetID</div>
                    <input type="text" placeholder="Like c5tuk9b765s2jfnbtfog" />
                </div>
                <div className="modalInput">
                    <div className="modalInputTitle">ListenPort</div>
                    <input type="number" placeholder="Like 3389" />
                </div>
                <div className="modalButton">
                    <div className="modalCancel" onClick={() => props.setModal(false)}>Cancel</div>
                    <div className="modalOK">Create</div>
                </div>
            </div>
        </div>
    );
}

export default ProxyModal;