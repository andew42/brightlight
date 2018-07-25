import * as React from "react";
import './ButtonPad.css';
import {Bouncefix} from "./Bouncefix";
import Button from "./Button";
import {Icon, Message} from "semantic-ui-react";

// Display an array of buttons and handle tap and press
export default class ButtonPad extends React.Component {

    render() {
        if (this.props.allButtons === undefined)
            return this.waitingForServer();

        return <Bouncefix className="Bouncefix">
            <div className="button-pad">{this.props.allButtons.map(button =>
                <Button key={button.key}
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
