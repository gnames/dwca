package dwca

func (a *arch) saveMetaOutput() error {
	bs, err := a.outputMeta.Bytes()
	if err != nil {
		return err
	}
	return a.dcFile.SaveToFile("meta.xml", bs)
}

func (a *arch) saveEmlOutput() error {
	bs, err := a.emlData.Bytes()
	if err != nil {
		return err
	}
	return a.dcFile.SaveToFile("eml.xml", bs)
}
