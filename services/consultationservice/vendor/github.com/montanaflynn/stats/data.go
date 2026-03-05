package stats


type Float64Data []float64


func (f Float64Data) Get(i int) float64 { return f[i] }


func (f Float64Data) Len() int { return len(f) }


func (f Float64Data) Less(i, j int) bool { return f[i] < f[j] }


func (f Float64Data) Swap(i, j int) { f[i], f[j] = f[j], f[i] }


func (f Float64Data) Min() (float64, error) { return Min(f) }


func (f Float64Data) Max() (float64, error) { return Max(f) }


func (f Float64Data) Sum() (float64, error) { return Sum(f) }


func (f Float64Data) CumulativeSum() ([]float64, error) { return CumulativeSum(f) }


func (f Float64Data) Mean() (float64, error) { return Mean(f) }


func (f Float64Data) Median() (float64, error) { return Median(f) }


func (f Float64Data) Mode() ([]float64, error) { return Mode(f) }


func (f Float64Data) GeometricMean() (float64, error) { return GeometricMean(f) }


func (f Float64Data) HarmonicMean() (float64, error) { return HarmonicMean(f) }


func (f Float64Data) MedianAbsoluteDeviation() (float64, error) {
	return MedianAbsoluteDeviation(f)
}


func (f Float64Data) MedianAbsoluteDeviationPopulation() (float64, error) {
	return MedianAbsoluteDeviationPopulation(f)
}


func (f Float64Data) StandardDeviation() (float64, error) {
	return StandardDeviation(f)
}


func (f Float64Data) StandardDeviationPopulation() (float64, error) {
	return StandardDeviationPopulation(f)
}


func (f Float64Data) StandardDeviationSample() (float64, error) {
	return StandardDeviationSample(f)
}


func (f Float64Data) QuartileOutliers() (Outliers, error) {
	return QuartileOutliers(f)
}


func (f Float64Data) Percentile(p float64) (float64, error) {
	return Percentile(f, p)
}


func (f Float64Data) PercentileNearestRank(p float64) (float64, error) {
	return PercentileNearestRank(f, p)
}


func (f Float64Data) Correlation(d Float64Data) (float64, error) {
	return Correlation(f, d)
}


func (f Float64Data) AutoCorrelation(lags int) (float64, error) {
	return AutoCorrelation(f, lags)
}


func (f Float64Data) Pearson(d Float64Data) (float64, error) {
	return Pearson(f, d)
}


func (f Float64Data) Quartile(d Float64Data) (Quartiles, error) {
	return Quartile(d)
}


func (f Float64Data) InterQuartileRange() (float64, error) {
	return InterQuartileRange(f)
}


func (f Float64Data) Midhinge(d Float64Data) (float64, error) {
	return Midhinge(d)
}


func (f Float64Data) Trimean(d Float64Data) (float64, error) {
	return Trimean(d)
}


func (f Float64Data) Sample(n int, r bool) ([]float64, error) {
	return Sample(f, n, r)
}


func (f Float64Data) Variance() (float64, error) {
	return Variance(f)
}


func (f Float64Data) PopulationVariance() (float64, error) {
	return PopulationVariance(f)
}


func (f Float64Data) SampleVariance() (float64, error) {
	return SampleVariance(f)
}


func (f Float64Data) Covariance(d Float64Data) (float64, error) {
	return Covariance(f, d)
}


func (f Float64Data) CovariancePopulation(d Float64Data) (float64, error) {
	return CovariancePopulation(f, d)
}


func (f Float64Data) Sigmoid() ([]float64, error) {
	return Sigmoid(f)
}



func (f Float64Data) SoftMax() ([]float64, error) {
	return SoftMax(f)
}


func (f Float64Data) Entropy() (float64, error) {
	return Entropy(f)
}


func (f Float64Data) Quartiles() (Quartiles, error) {
	return Quartile(f)
}
