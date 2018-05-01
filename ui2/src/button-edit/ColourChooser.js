import * as React from "react";
import './ColourChooser.css';
import {Button, Modal} from "semantic-ui-react";
import "../colour/ColourConversion.js"
import VerticalSlider from "../buttons/VerticalSlider";
import {asColourObject, asColourString, HSVtoRGB, RGBtoHSV} from "../colour/ColourConversion";
import MountNotifier from "./MountNotifier";

export default class ColourChooser extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            hue: 0,
            saturation: 0,
            value: 0
        };
    }

    setInitialColour() {
        let initialColour = RGBtoHSV(asColourObject(this.props.colour));
        this.setState({hue: initialColour.h, saturation: initialColour.s, value: initialColour.v})
    }

    render() {
        let colour = asColourString(HSVtoRGB(this.state.hue, this.state.saturation, this.state.value));
        let saturatedColour = asColourString(HSVtoRGB(this.state.hue, 100, 100));
        return <Modal trigger={this.props.trigger}
                      open={this.props.isOpen}>
            <Modal.Header>Choose Colour</Modal.Header>
            <Modal.Content>
                <MountNotifier componentDidMount={() => this.setInitialColour()}/>
                <Modal.Description>
                    <div className='cp-content-height'>
                        <div className='cp-sliders'>
                            <VerticalSlider className='cp-slider cp-slider-hue'
                                            min={0}
                                            max={360}
                                            pos={this.state.hue}
                                            onPosChange={p => this.setState({hue: p})}/>

                            <VerticalSlider className='cp-slider cp-slider-saturation'
                                            min={0}
                                            max={100}
                                            pos={this.state.saturation}
                                            onPosChange={p => this.setState({saturation: p})}
                                            sliderColour={saturatedColour}/>

                            <VerticalSlider className='cp-slider cp-slider-value'
                                            min={0}
                                            max={100}
                                            pos={this.state.value}
                                            onPosChange={p => this.setState({value: p})}
                                            sliderColour={saturatedColour}/>
                        </div>
                    </div>
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions style={{background: colour}}>
                <Button primary content='OK' onClick={() => console.info('OK')}/>
                <Button content='Cancel' onClick={this.props.onCancel}/>
            </Modal.Actions>
        </Modal>
    }
}
