import * as React from "react";
import * as Hammer from "hammerjs"

export default class Button extends React.Component {
    constructor(props) {
        super(props);

        let emptyFunction = function () {
        };

        // Wire up onTap handler if provided
        if (props.onTap !== undefined) {
            this.onTap = props.onTap;
        } else {
            this.onTap = emptyFunction;
        }
        this.onTap = this.onTap.bind(this);

        // Wire up onPress handler if provided
        if (props.onPress !== undefined) {
            this.onPress = props.onPress;
        } else {
            this.onPress = emptyFunction;
        }
        this.onPress = this.onPress.bind(this);
    }

    componentDidMount() {
        this.hammer = Hammer(this.domButton);
        this.hammer.get('press').set({time: 1000});
        this.hammer.on('tap', this.onTap);
        this.hammer.on('press', this.onPress);
    }

    componentWillUnmount() {
        this.hammer.off('tap', this.onTap);
        this.hammer.off('press', this.onPress);
    }

    render() {
        return <button ref={(domElement) => {
            this.domButton = domElement;
        }}>
            {this.props.label}
        </button>
    }
}
