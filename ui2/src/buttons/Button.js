import * as React from "react";
import './Button.css';

export default class Button extends React.Component {
    constructor(props) {
        super(props);

        this.state = {pressed: false};
    }

    onTouchStart() {
        this.timerId = setTimeout(() => this.setState({pressed: true}), 1000);
    }

    onTouchEnd(e) {
        clearTimeout(this.timerId);
        if (this.state.pressed) {
            e.preventDefault();
            this.setState({pressed: false});
            this.props.onPressUp();
        }
        else {
            this.props.onTap();
        }
    }

    render() {
        return <button style={this.props.style}
                       onTouchStart={() => this.onTouchStart()}
                       onTouchEnd={e => this.onTouchEnd(e)}
                       onMouseDown={() => this.onTouchStart()}
                       onMouseUp={e => this.onTouchEnd(e)}
                       className={this.state.pressed ? 'button-long-pressed' : 'button-pressed'}
        >
            <div>{this.props.label}</div>
        </button>
    }
}
