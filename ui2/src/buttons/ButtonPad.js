import * as React from "react";
import './ButtonPad.css';
import {Bouncefix} from "./Bouncefix";
import Button from "./Button";
import {Icon, Message} from "semantic-ui-react";
import {OpenWebSocket} from "../server-proxy/webSocket";

// Display an array of buttons and handle tap and press
export default class ButtonPad extends React.Component {
    constructor(props) {
        super(props);
        this.state = {activeButtonKey: 0};
    }

    ws = null;

    // Open web socket to stream active button state transitions
    componentDidMount() {
        this.ws = OpenWebSocket('ButtonState', b => {
            console.info('Active button changed to ' + b);
            this.setState({activeButtonKey: b});
        });
    }

    // Close web socket
    componentWillUnmount() {
        this.ws.close()
    }

    render() {
        if (this.props.allButtons === undefined)
            return this.waitingForServer();

        return <Bouncefix className="Bouncefix">
            <div className="button-pad">{this.props.allButtons.map(button =>
                <Button key={button.key}
                        active={button.key === this.state.activeButtonKey}
                        onTap={() => this.props.onButtonTap(button.key)}
                        onPressUp={() => this.props.onButtonPress(this.props.history, button.key)}
                        label={button.name}/>)}
            </div>
            <div className="button-pad-bottom-padding"/>
        </Bouncefix>
    }

    waitingForServer() {
        return <Message warning icon>
            <Icon name='circle notched' loading/>
            <Message.Content>
                <Message.Header>Waiting for Server</Message.Header>
                Ensure phone is connected to correct WiFi
            </Message.Content>
        </Message>
    }
}
