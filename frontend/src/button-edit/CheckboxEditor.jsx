import * as React from "react";
import './CheckboxEditor.css';

export function CheckboxEditor(props) {
    return (
        <label className='cbe-label'>
            <input type='checkbox'
                   className='cbe-checkbox'
                   checked={props.checked}
                   onChange={e => props.onChange(e.target.checked)}/>
            {props.label}
        </label>
    );
}
