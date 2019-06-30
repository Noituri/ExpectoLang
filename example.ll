declare float @printf(float %x)

define float @main() {
entry:
  %calltmp = call float @printf(float 3.000000e+00)
  ret float 5.000000e+00
}
