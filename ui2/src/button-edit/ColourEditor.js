import * as React from "react";
import './ColourEditor.css';
import {Button, Label} from "semantic-ui-react";
import ColourChooser from "./ColourChooser";
import PropTypes from 'prop-types';

export default class ColourEditor extends React.Component {

    static propTypes = {
        // current colour
        colour: PropTypes.oneOfType([PropTypes.number, PropTypes.string, PropTypes.object]).isRequired,

        // function to call when colour is changed
        onChangeColour: PropTypes.func.isRequired,
    };

    render() {
        return <div className='colour-editor-container'>
            <ColourChooser trigger={<Button style={{backgroundColor: this.props.colour}}/>}
                           colour={this.props.colour}
                           onOk={colour => {
                               this.props.onChangeColour(colour);
                           }}/>
            <Label pointing='left' content='Colour'/>
        </div>
    }
}
