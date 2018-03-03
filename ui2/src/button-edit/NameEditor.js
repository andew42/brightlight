import * as React from "react";

export default class NameEditor extends React.Component {
    constructor(props) {
        super(props);
        this.button = props.button;
    }

    render() {
        return <div>
            <label>
                Name
                <input type="text" value={this.props.name} onChange={e => this.props.onNameChanged(e.target.value)}/>
            </label>
        </div>
    }
}
