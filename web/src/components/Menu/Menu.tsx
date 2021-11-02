import React from 'react';
import { BsFillMenuAppFill } from "react-icons/bs";
import { MdSettingsPower } from "react-icons/md";
import "./index.css";

const Menu: React.FunctionComponent = (): JSX.Element => {
    return <div className="menu">
        <div className="Logo">
            <img src="/images/logo.png" alt="" />
        </div>
        <div className="Item selected"><BsFillMenuAppFill /></div>
        <div className="Item logout"><MdSettingsPower /></div>
    </div>
}

export default Menu;