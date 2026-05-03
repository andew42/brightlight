import * as React from "react";
import './NameEditor.css';

export function NameEditor(props) {
    return (
        <fieldset className='ne-section'>
            <legend>Name</legend>
            <input className='ne-input'
                   value={props.name}
                   onChange={e => props.onNameChanged(e.target.value)}/>
            {props.error && <div className='ne-error'>{props.error}</div>}
        </fieldset>
    );
}
