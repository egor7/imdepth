#version 3.7;
global_settings {
  assumed_gamma 0.5
}

#include "colors.inc"

//light_source{ < 2000,-1000,-1000> White }
//light_source{ <-2000,-1000,-1000> White }
//light_source{ <    0,-1000, 2000> White }

//light_source{ <    0,-1000, 0> White }
background { color White }

camera {
  location  <-5,10,-8>
  look_at   <0,-0.5,0>
  angle 30
}

height_field {
    png "heights-b.png"
    smooth
    pigment {
        image_map{png "colors.png"}
        rotate x*90
    }
    translate <-.5, -.5, -.5>
    scale <5, 0.2, 5>
  }



