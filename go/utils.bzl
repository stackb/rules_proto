
def get_importmappings(mappings):
    """Build list of importmappings

    Args:
        mappings: dict<string,string> src: dst 
    Returns
        list<string> in the form 'M{src}={dst}`
    """
    if not mappings:
        return []
    return ["M%s=%s" % (src, dst) for src, dst in mappings.items()]
