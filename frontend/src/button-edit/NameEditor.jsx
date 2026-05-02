import * as React from "react";
import {Input, Label, Segment} from "semantic-ui-react";
import './NameEditor.css';

export function NameEditor(props) {
    return <Segment color='blue' attached style={{textAlign: 'left'}}>
        <Label color='blue' attached='top left' size='large' content='Name'/>
        <Input fluid value={props.name}
               onChange={e => props.onNameChanged(e.target.value)}/>
        {props.error !== undefined && <Label pointing basic color='red' content={props.error}/>}
    </Segment>
}
