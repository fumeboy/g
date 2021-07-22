package util

import "github.com/dave/dst/decorator"

func DecCopy(dec_src,dec_dst *decorator.Decorator){
	for k,v := range dec_src.Ast.Nodes{
		dec_dst.Ast.Nodes[k] = v
	}
	for k,v := range dec_src.Ast.Scopes{
		dec_dst.Ast.Scopes[k] = v
	}
	for k,v := range dec_src.Ast.Objects{
		dec_dst.Ast.Objects[k] = v
	}
	for k,v := range dec_src.Dst.Nodes{
		dec_dst.Dst.Nodes[k] = v
	}
	for k,v := range dec_src.Dst.Scopes{
		dec_dst.Dst.Scopes[k] = v
	}
	for k,v := range dec_src.Dst.Objects{
		dec_dst.Dst.Objects[k] = v
	}
}
