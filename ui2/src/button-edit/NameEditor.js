import * as React from "react";
import {Input, Label, Segment} from "semantic-ui-react";

export function NameEditor(props) {
    return <Segment color='blue' attached>
        <Label color='blue' attached='top left'>
            Name
        </Label>
        <Input fluid value={props.name}
               onChange={e => props.onNameChanged(e.target.value)}/>
    </Segment>
}
