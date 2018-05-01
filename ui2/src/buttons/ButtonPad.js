import * as React from "react";
import './ButtonPad.css';
import {getButtons} from "../server-proxy/buttons";
import {runAnimation} from "../server-proxy/animation";
import {Bouncefix} from "./Bouncefix";
import Button from "./Button";

// Display an array of buttons retrieved from the server
export default class ButtonPad extends React.Component {
    constructor(props) {
        super(props);
        this.state = {buttons: []};
        this.history = props.history
    }

    componentDidMount() {
        // Retrieve button state from server to update our state
        getButtons((buttons) => this.setState({buttons: buttons}),
            (xhr) => console.error(xhr))
    }

    render() {
        return <Bouncefix className="Bouncefix">
            <div>
                <div className="button-pad">{this.state.buttons.map((button) =>
                    <Button key={button.name}
                            onTap={() => runAnimation(button.segments)}
                            onPress={() => {
                                this.history.push('/button-edit', {'button': button, 'buttons': this.state.buttons})
                            }}
                            label={button.name}/>)}
                </div>
            </div>
        </Bouncefix>
    }
}

