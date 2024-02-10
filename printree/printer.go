package printree

func printRow(
	st IPrintableTreeNode,
	currPrefix string,
	seqInChildren int,
	nodeSerializer PrintableNodeSerializer,
	setting *SerializeSetting,
) (canContinue bool, err error) {
	if _, err = setting.Println(
		nodeSerializer.mkLine(st, currPrefix, seqInChildren),
	); err != nil {
		return false, err
	}

	if setting.DrillDown != nil && setting.Predictor != nil {
		if !setting.Predictor(st) {
			if setting.DrillDown.ShowMark != "" && st.ChildrenCount() > 0 {
				if _, err = setting.Println(
					nodeSerializer.mkMark(st, currPrefix, setting.DrillDown.ShowMark),
				); err != nil {
					return false, err
				}
			}
			return false, nil
		}
	}
	return true, nil
}

func recursivePrint(
	st IPrintableTreeNode,
	prefix string,
	seqInChildren, childrenCount int,
	nodeSerializer PrintableNodeSerializer,
	setting *SerializeSetting,
) error {
	currPrefix, nextLevelPrefix := setting.TreeMarks.BuildPrefixes(prefix, seqInChildren == childrenCount-1)
	canContinue, err := printRow(st, currPrefix, seqInChildren, nodeSerializer, setting)
	if err != nil {
		return err
	}
	if !canContinue {
		return nil
	}

	if err = st.ForEachPrintableChild(func(
		child IPrintableTreeNode,
		ind int,
	) error {
		return recursivePrint(child, nextLevelPrefix, ind, st.ChildrenCount(), nodeSerializer, setting)
	}); err != nil {
		return err
	}

	return nil
}
