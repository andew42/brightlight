import * as React from "react";
import {Button, Label} from "semantic-ui-react";
import './ColourEditor.css';
import ColourChooser from "./ColourChooser";

export default class ColourEditor extends React.Component {
    // props: colour, onColourChanged
    render() {
        return <div className='colour-editor-container'>
            <ColourChooser trigger={<Button className='colour-editor-button'
                                            style={{backgroundColor: this.props.colour.asColourString()}}/>}
                           colour={this.props.colour}
                           onColourChanged={c => this.props.onColourChanged(c)}/>
            <Label pointing='left' content='Colour'/>
        </div>
    }
}
