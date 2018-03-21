import * as React from "react";
import {Input, Label, Segment} from "semantic-ui-react";

// Edit a name property
export default class NameEditor extends React.Component {
    render() {
        return <Segment color='blue' attached>
            <Label color='blue' attached='top left'>Name</Label>
            <Input fluid value={this.props.name}
                   onChange={e => this.props.onNameChanged(e.target.value)}/>
        </Segment>
    }
}
