import * as React from "react";
import './ColourEditor.css';
import {Label} from "semantic-ui-react";

export function ColourEditor(props) {

    return <div className='colour-editor-container'>
        <div style={{backgroundColor: props.colour}} className='colour-patch'/>
        <Label pointing='left' content='Colour'/>
    </div>
}
