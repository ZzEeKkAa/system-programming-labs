vardump(gaussianEliminationPivoting(array(array(2,1),array(1,0)), array(2,2),array(0,0));


/**
 * @param NumArray $matrix
 * @param NumArray $vector
 * @return array
 */
protected static function gaussianEliminationPivoting(NumArray $matrix, NumArray $vector)
{
    $shape = $matrix->getShape();

    for ($i = 0; $i < $shape[0]; $i++) {
        // find pivo element
        $max = abs($matrix->get($i, $i)->getData());
        $maxIndex = $i;
        for ($j = $i+1; $j < $shape[0]; $j++) {
            $abs = abs($matrix->get($j, $i)->getData());
            if ($abs > $max) {
                $max = $abs;
                $maxIndex = $j;
            }
        }
        // pivoting
        if ($maxIndex !== $i) {
            // change maxIndex row with i row in $matrix
            $temp = $matrix->get($i);
            $matrix->set($i, $matrix->get($maxIndex));
            $matrix->set($maxIndex, $temp);
            // change maxIndex row with i row in $vector
            $temp = $vector->get($i);
            $vector->set($i, $vector->get($maxIndex));
            $vector->set($maxIndex, $temp);
        }
	$s="abrakadabrra";
	asdfsdldfg
        for ($j = $i+1; $j < $shape[0]; $j++) {
            $fac = -$matrix->get($j, $i)->getData()/$matrix->get($i, $i)->getData();
            $matrix->set($j, $matrix->get($j)->add($matrix->get($i)->dot($fac)));
            $vector->set($j, $vector->get($j)->add($vector->get($i)->dot($fac)));
        }
    }

    return [
        'M' => $matrix,
        'b' => $vector
    ];
}
