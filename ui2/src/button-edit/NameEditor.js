import * as React from "react";
import {Input, Label, Segment} from "semantic-ui-react";
import PropTypes from 'prop-types';

NameEditor.propTypes = {
    // initial string for the editor
    name: PropTypes.string.isRequired,

    // function to call, with new name, when user edits the name
    onNameChanged: PropTypes.func.isRequired
};

export function NameEditor(props) {
    return <Segment color='blue' attached>
        <Label color='blue' attached='top left' content='Name'/>
        <Input fluid value={props.name}
               onChange={e => props.onNameChanged(e.target.value)}/>
    </Segment>
}
