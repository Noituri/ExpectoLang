; ModuleID = 'expectoroot'
source_filename = "expectoroot"

declare float @sin(float %x)

define float @main() {
entry:
  %calltmp = call float @sin(float 3.000000e+00)
  ret float 5.000000e+00
}
