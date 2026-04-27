import * as React from "react";

export default class Button extends React.Component {
    constructor(props) {
        super(props);

        this.state = {pressed: false};
    }

    onTouchStart() {
        this.timerId = setTimeout(() => this.setState({pressed: true}), 1000);
        // Send the tap immediately as it doesn't cause significant UI change
        this.props.onTap();
    }

    onTouchEnd(e) {
        clearTimeout(this.timerId);
        if (this.state.pressed) {
            e.preventDefault();
            this.setState({pressed: false});
            // Send the press on touch end because this can cause UI navigation
            // and on mobile safari if we send the press immediately after the
            // timer expires we see the touch end over the new screen and this
            // is sometimes actioned which is very confusing immediately after
            // screen navigation
            this.props.onPressUp();
        }
    }

    render() {
        return <button onTouchStart={() => this.onTouchStart()}
                       onTouchEnd={e => this.onTouchEnd(e)}
                       onMouseDown={() => this.onTouchStart()}
                       onMouseUp={e => this.onTouchEnd(e)}
                       className={this.state.pressed ? 'button-pressed' : ''}>
            <div>{this.props.label}</div>
        </button>
    }
}
