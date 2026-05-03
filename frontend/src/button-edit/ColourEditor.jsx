import * as React from "react";
import './ColourEditor.css';
import ColourChooser from "./ColourChooser";

export default class ColourEditor extends React.Component {
    // props: colour, label, onColourChanged
    render() {
        return <div className='colour-editor-container'>
            <ColourChooser trigger={<button className='colour-editor-button'
                                            style={{backgroundColor: this.props.colour.asColourString()}}/>}
                           colour={this.props.colour}
                           onColourChanged={c => this.props.onColourChanged(c)}/>
            <span className='colour-editor-label'>{this.props.label}</span>
        </div>;
    }
}
