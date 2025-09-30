<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Pricing\V1;

/**
 */
class PricingServiceClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Pricing\V1\QuoteRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function Quote(\Pricing\V1\QuoteRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/pricing.v1.PricingService/Quote',
        $argument,
        ['\Pricing\V1\QuoteResponse', 'decode'],
        $metadata, $options);
    }

}
