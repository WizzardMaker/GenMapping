# Tags
AutoMapper can be configured with tags. Each tag corresponds to a different operation or configuration

## Mapping Method Tags
> ## Translation (@translate)
> The translation tag directly maps a given source parameter to the target field
> 
> _Beware - this operation is not type checked by AutoMapper and could result in incorrect/not compileable output_
> ### Syntax
> ```go
> //@translate(from="A", to="B")
> ```
> _from_ - Source field, must be a parameter of the mapping function.
> 
> _to_ - Destination field of the output object

> ## Expression (@expression)
> The expression tag allows to define custom code for the mapping of the target field
> 
> _Beware - this operation is not checked by AutoMapper and could result in incorrect/not compileable output_
> ### Syntax
> ```
> //@expression(expression="", target="", isType="false|true")
> ```
> _expression_ - the desired Go expression
> 
> _target_ - the field name to which to apply that expression