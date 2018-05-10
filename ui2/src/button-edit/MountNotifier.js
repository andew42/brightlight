import * as React from "react";
import PropTypes from 'prop-types';

export default class MountNotifier extends React.Component {

    static propTypes = {
        // function to call when component mounts
        componentDidMount: PropTypes.func.isRequired,
    };

    // used to determine when a modal opens
    componentDidMount() {
        this.props.componentDidMount();
    }

    render() {
        // renter nothing
        return '';
    }
}
