; ModuleID = 'root'
source_filename = "root"

declare float @printf(float %x)

define float @main() {
entry:
	  %calltmp = call float @printf(float 4.000000e+00)
	    ret float %calltmp
}

