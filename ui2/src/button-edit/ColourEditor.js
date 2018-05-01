import * as React from "react";
import './ColourEditor.css';
import {Label} from "semantic-ui-react";
import Button from "../buttons/Button";
import ColourChooser from "./ColourChooser";

export default class ColourEditor extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isChooserOpen: false
        };
    }

    render() {
        return <div className='colour-editor-container'>
            <ColourChooser
                trigger={
                    <Button onTap={() => this.setState({isChooserOpen: true})}
                            style={{backgroundColor: this.props.colour}}/>
                }
                isOpen={this.state.isChooserOpen}
                colour={this.props.colour}
                onOk={colour => {
                    this.props.onChangeColour(colour);
                    this.setState({isChooserOpen: false});
                }}
                onCancel={() => this.setState({isChooserOpen: false})}/>
            <Label pointing='left' content='Colour'/>
        </div>
    }
}
