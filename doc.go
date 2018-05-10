package boneful

/*

The boneful package exists to provide a sort of impedance matching between
the simple and fast bone http multiplexer
(see https://godoc.org/github.com/go-zoo/bone)
and the complex and somewhat overbuilt go-restful package
(see https://github.com/emicklei/go-restful).

The problem is that go-restful imposes a lot of requirements on its
clients -- the worst of them is that it doesn't use go-standard
http handlers.

However, it DOES provide a nice framework for defining and documenting
restful endpoints. It just also requires the installation of go-swagger
and all of its support -- and then go-swagger doesn't play nice with
microservices. Consequently, we stopped using go-swagger very early,
and stopped using go-restful a while later.

However, we still wanted good documentation. So this package was developed
to provide a way to define a restful API that is almost completely
compatible with go-restful, but actually generates the setup for
a bone muxer.

This makes it easy to adapt an existing bone-based API to it simply by
changing the definitions of the endpoints in a way that also creates
the documentation for those endpoints.

It also makes it fairly easy to port go-restful-based APIs by minor mods
to the API definitions and somewhat larger mods to the handlers.

Finally, it automatically adds a /md endpoint to the API which will
deliver the documentation in .md (markdown) form, intended to be compatible
with GitHub's GFM.

*/
