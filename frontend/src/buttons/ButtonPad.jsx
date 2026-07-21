import * as React from "react";
import './ButtonPad.css';
import Button from "./Button";
import NoScroll from "./NoScroll";

// Display an array of buttons and handle tap and press
export default class ButtonPad extends React.Component {

    constructor(props) {
        super(props);
        this.dom = React.createRef();
    }

    render() {
        if (this.props.allButtons === undefined)
            return ButtonPad.noServer();

        return <NoScroll>
            <div className="button-pad">
                {this.props.allButtons.map(button =>
                    <div key={button.key}>
                        <Button active={button.key === this.props.activeButtonKey}
                                onTap={() => this.props.onButtonTap(button.key)}
                                onPressUp={() => this.props.onButtonPress(button.key)}
                                label={button.name}/>
                        <div className={button.key === this.props.activeButtonKey ? 'active' : ''}/>
                    </div>)}
            </div>
        </NoScroll>;
    }

    static noServer() {
        return (
            <div className='bp-no-server'>
                <div className='bp-spinner'/>
                <div className='bp-no-server-title'>Server Not Found</div>
                <div>Ensure phone is connected to WiFi and try again</div>
            </div>
        );
    }
}
