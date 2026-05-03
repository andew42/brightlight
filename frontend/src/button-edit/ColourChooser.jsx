import * as React from "react";
import './ColourChooser.css';
import Dialog from "../dialog/Dialog";
import VerticalSlider from "../buttons/VerticalSlider";
import MountNotifier from "./MountNotifier";
import Colour from "../colour/Colour";

export default class ColourChooser extends React.Component {
    // props: colour, trigger, onColourChanged

    constructor(props) {
        super(props);
        this.state = {open: false};
    }

    render() {
        let colour = new Colour(this.props.colour);
        let saturatedColour = new Colour({h: colour.h, s: 100, v: 100});
        return <>
            <span onClick={() => this.setState({open: true})}>{this.props.trigger}</span>
            <Dialog open={this.state.open}
                    header='Choose Colour'
                    actions={[
                        {key: 'ok', content: 'OK', primary: true, onClick: () => this.setState({open: false})},
                        {key: 'cancel', content: 'Cancel', onClick: () => {
                            this.props.onColourChanged(this.initialColour);
                            this.setState({open: false});
                        }}
                    ]}
                    onClose={() => this.setState({open: false})}>
                <MountNotifier componentDidMount={() => this.initialColour = colour}/>
                <div className='cp-content-height'>
                    <ColourSwatch colour={colour.asColourString()}/>
                    <div className='cp-sliders'>
                        <VerticalSlider className='cp-slider cp-slider-hue'
                                        min={0} max={360} pos={colour.h}
                                        onPosChange={p => this.props.onColourChanged(
                                            new Colour({h: p, s: colour.s, v: colour.v}))}/>
                        <VerticalSlider className='cp-slider cp-slider-saturation'
                                        min={0} max={100} pos={colour.s}
                                        onPosChange={p => this.props.onColourChanged(
                                            new Colour({h: colour.h, s: p, v: colour.v}))}
                                        sliderColour={saturatedColour.asColourString()}/>
                        <VerticalSlider className='cp-slider cp-slider-value'
                                        min={0} max={100} pos={colour.v}
                                        onPosChange={p => this.props.onColourChanged(
                                            new Colour({h: colour.h, s: colour.s, v: p}))}
                                        sliderColour={saturatedColour.asColourString()}/>
                    </div>
                </div>
            </Dialog>
        </>;
    }
}

function ColourSwatch({colour}) {
    return <div className='cp-colour-swatch' style={{backgroundColor: colour}}/>;
}
