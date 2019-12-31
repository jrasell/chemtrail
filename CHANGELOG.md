## 0.1.0 (Unreleased)

IMPROVEMENTS:
 * Add systems metrics API function and CLI [[GH-16]](https://github.com/jrasell/chemtrail/pull/16)
 * Add NoOp client provider to allow scaling evaluations where results are logged and cluster state is not altered [[GH-18]](https://github.com/jrasell/chemtrail/pull/18)

BUG FIXES:
 * Do not log Chemtrail allocation nodeID if Chemtrail is not found to be running on Nomad [[GH-23]](https://github.com/jrasell/chemtrail/pull/23)
 * Do not log AWS ASG provider setup when it is not enabled [[GH-20]](https://github.com/jrasell/chemtrail/pull/20)
 * Correctly format ProviderConfig CLI output when reading a scaling policy [[GH-22]](https://github.com/jrasell/chemtrail/pull/22)

## 0.0.1 (17 December, 2019)

 * Initial release
