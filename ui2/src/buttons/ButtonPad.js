import * as React from "react";
import './ButtonPad.css';
import {Bouncefix} from "./Bouncefix";
import Button from "./Button";

// Display an array of buttons and handle tap and press
export default class ButtonPad extends React.Component {

    render() {
        return <Bouncefix className="Bouncefix">
            <div>
                <div className="button-pad">{this.props.allButtons.map(button =>
                    <Button key={button.key}
                            onTap={() => this.props.onButtonTap(button.key)}
                            onPress={() => this.props.onButtonPress(this.props.history, button.key)}
                            label={button.name}/>)}
                </div>
            </div>
        </Bouncefix>
    }
}
